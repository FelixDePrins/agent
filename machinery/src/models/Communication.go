package models

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/kerberos-io/joy4/av/pubsub"
	"github.com/kerberos-io/joy4/cgo/ffmpeg"
	"github.com/tevino/abool"
)

// The communication struct that is managing
// all the communication between the different goroutines.
type Communication struct {
	Context               *context.Context
	CancelContext         *context.CancelFunc
	PackageCounter        *atomic.Value
	LastPacketTimer       *atomic.Value
	CloudTimestamp        *atomic.Value
	HandleBootstrap       chan string
	HandleStream          chan string
	HandleSubStream       chan string
	HandleMotion          chan MotionDataPartial
	HandleAudio           chan AudioDataPartial
	HandleUpload          chan string
	HandleHeartBeat       chan string
	HandleLiveSD          chan int64
	HandleLiveHDKeepalive chan string
	HandleLiveHDHandshake chan RequestHDStreamPayload
	HandleLiveHDPeers     chan string
	HandleONVIF           chan OnvifAction
	IsConfiguring         *abool.AtomicBool
	Queue                 *pubsub.Queue
	SubQueue              *pubsub.Queue
	DecoderMutex          *sync.Mutex
	SubDecoderMutex       *sync.Mutex
	Decoder               *ffmpeg.VideoDecoder
	SubDecoder            *ffmpeg.VideoDecoder
	Image                 string
	CameraConnected       bool
	HasBackChannel        bool
}
