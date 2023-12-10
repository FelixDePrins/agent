package http

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/kerberos-io/agent/machinery/src/capture"
	"github.com/kerberos-io/agent/machinery/src/onvif"
	"github.com/kerberos-io/agent/machinery/src/routers/websocket"

	"github.com/kerberos-io/agent/machinery/src/cloud"
	configService "github.com/kerberos-io/agent/machinery/src/config"
	"github.com/kerberos-io/agent/machinery/src/models"
	"github.com/kerberos-io/agent/machinery/src/utils"
)

func AddRoutes(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware, configDirectory string, configuration *models.Configuration, communication *models.Communication, captureDevice *capture.Capture) *gin.RouterGroup {

	r.GET("/ws", func(c *gin.Context) {
		websocket.WebsocketHandler(c, communication, captureDevice)
	})

	// This is legacy should be removed in future! Now everything
	// lives under the /api prefix.
	r.GET("/config", func(c *gin.Context) {

		// We'll try to get a snapshot from the camera.
		base64Image := capture.Base64Image(captureDevice, communication)
		if base64Image != "" {
			communication.Image = base64Image
		}

		c.JSON(200, gin.H{
			"config":   configuration.Config,
			"custom":   configuration.CustomConfig,
			"global":   configuration.GlobalConfig,
			"snapshot": communication.Image,
		})
	})

	// This is legacy should be removed in future! Now everything
	// lives under the /api prefix.
	r.POST("/config", func(c *gin.Context) {
		var config models.Config
		err := c.BindJSON(&config)
		if err == nil {
			err := configService.SaveConfig(configDirectory, config, configuration, communication)
			if err == nil {
				c.JSON(200, gin.H{
					"data": "☄ Reconfiguring",
				})
			} else {
				c.JSON(400, gin.H{
					"data": "Something went wrong: " + err.Error(),
				})
			}
		} else {
			c.JSON(400, gin.H{
				"data": "Something went wrong: " + err.Error(),
			})
		}
	})

	api := r.Group("/api")
	{
		api.POST("/login", authMiddleware.LoginHandler)

		api.GET("/dashboard", func(c *gin.Context) {

			// Check if camera is online.
			cameraIsOnline := communication.CameraConnected

			// If an agent is properly setup with Kerberos Hub, we will send
			// a ping to Kerberos Hub every 15seconds. On receiving a positive response
			// it will update the CloudTimestamp value.
			cloudIsOnline := false
			if communication.CloudTimestamp != nil && communication.CloudTimestamp.Load() != nil {
				timestamp := communication.CloudTimestamp.Load().(int64)
				if timestamp > 0 {
					cloudIsOnline = true
				}
			}

			// The total number of recordings stored in the directory.
			recordingDirectory := configDirectory + "/data/recordings"
			numberOfRecordings := utils.NumberOfMP4sInDirectory(recordingDirectory)

			// All days stored in this agent.
			days := []string{}
			latestEvents := []models.Media{}
			files, err := utils.ReadDirectory(recordingDirectory)
			if err == nil {
				events := utils.GetSortedDirectory(files)

				// Get All days
				days = utils.GetDays(events, recordingDirectory, configuration)

				// Get all latest events
				var eventFilter models.EventFilter
				eventFilter.NumberOfElements = 5
				latestEvents = utils.GetMediaFormatted(events, recordingDirectory, configuration, eventFilter) // will get 5 latest recordings.
			}

			c.JSON(200, gin.H{
				"offlineMode":        configuration.Config.Offline,
				"cameraOnline":       cameraIsOnline,
				"cloudOnline":        cloudIsOnline,
				"numberOfRecordings": numberOfRecordings,
				"days":               days,
				"latestEvents":       latestEvents,
			})
		})

		api.POST("/latest-events", func(c *gin.Context) {
			var eventFilter models.EventFilter
			err := c.BindJSON(&eventFilter)
			if err == nil {
				// Default to 10 if no limit is set.
				if eventFilter.NumberOfElements == 0 {
					eventFilter.NumberOfElements = 10
				}
				recordingDirectory := configDirectory + "/data/recordings"
				files, err := utils.ReadDirectory(recordingDirectory)
				if err == nil {
					events := utils.GetSortedDirectory(files)
					// We will get all recordings from the directory (as defined by the filter).
					fileObjects := utils.GetMediaFormatted(events, recordingDirectory, configuration, eventFilter)
					c.JSON(200, gin.H{
						"events": fileObjects,
					})
				} else {
					c.JSON(400, gin.H{
						"data": "Something went wrong: " + err.Error(),
					})
				}
			} else {
				c.JSON(400, gin.H{
					"data": "Something went wrong: " + err.Error(),
				})
			}
		})

		api.GET("/days", func(c *gin.Context) {
			recordingDirectory := configDirectory + "/data/recordings"
			files, err := utils.ReadDirectory(recordingDirectory)
			if err == nil {
				events := utils.GetSortedDirectory(files)
				days := utils.GetDays(events, recordingDirectory, configuration)
				c.JSON(200, gin.H{
					"events": days,
				})
			} else {
				c.JSON(400, gin.H{
					"data": "Something went wrong: " + err.Error(),
				})
			}
		})

		api.GET("/config", func(c *gin.Context) {

			// We'll try to get a snapshot from the camera.
			base64Image := capture.Base64Image(captureDevice, communication)
			if base64Image != "" {
				communication.Image = base64Image
			}

			c.JSON(200, gin.H{
				"config":   configuration.Config,
				"custom":   configuration.CustomConfig,
				"global":   configuration.GlobalConfig,
				"snapshot": communication.Image,
			})
		})

		api.POST("/config", func(c *gin.Context) {
			var config models.Config
			err := c.BindJSON(&config)
			if err == nil {
				err := configService.SaveConfig(configDirectory, config, configuration, communication)
				if err == nil {
					c.JSON(200, gin.H{
						"data": "☄ Reconfiguring",
					})
				} else {
					c.JSON(200, gin.H{
						"data": "☄ Reconfiguring",
					})
				}
			} else {
				c.JSON(400, gin.H{
					"data": "Something went wrong: " + err.Error(),
				})
			}
		})

		api.GET("/restart", func(c *gin.Context) {
			communication.HandleBootstrap <- "restart"
			c.JSON(200, gin.H{
				"restarted": true,
			})
		})

		api.GET("/stop", func(c *gin.Context) {
			communication.HandleBootstrap <- "stop"
			c.JSON(200, gin.H{
				"stopped": true,
			})
		})

		api.POST("/onvif/verify", func(c *gin.Context) {
			onvif.VerifyOnvifConnection(c)
		})

		api.POST("/hub/verify", func(c *gin.Context) {
			cloud.VerifyHub(c)
		})

		api.POST("/persistence/verify", func(c *gin.Context) {
			cloud.VerifyPersistence(c, configDirectory)
		})

		// Camera specific methods. Doesn't require any authorization.
		// These are available for anyone, but require the agent, to reach
		// the camera.
		api.POST("/camera/onvif/login", LoginToOnvif)
		api.POST("/camera/onvif/capabilities", GetOnvifCapabilities)
		api.POST("/camera/onvif/presets", GetOnvifPresets)
		api.POST("/camera/onvif/gotopreset", GoToOnvifPreset)
		api.POST("/camera/onvif/pantilt", DoOnvifPanTilt)
		api.POST("/camera/onvif/zoom", DoOnvifZoom)
		api.POST("/camera/verify/:streamType", capture.VerifyCamera)

		// Secured endpoints..
		api.Use(authMiddleware.MiddlewareFunc())
		{
		}
	}
	return api
}
