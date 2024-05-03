import React from 'react';
import { DateTimePickerComponent } from '@syncfusion/ej2-react-calendars';
import './TimePicker.scss';
import { t } from 'i18next';
import PropTypes from 'prop-types';

class TimePicker extends React.PureComponent {
  maxDate = new Date();

  constructor(props) {
    super(props);
    this.state = {
      Date: new Date(),
      placeholderTranslate: t('timepicker.placeholder'),
    };
    this.handleChange = this.handleChange.bind(this);
  }

  handleChange = (event) => {
    console.log(event);
    if (event.value !== null) {
      const elements = document.getElementsByClassName(
        'stats grid-container --four-columns'
      );
      while (elements.length > 0) {
        elements[0].parentNode.removeChild(elements[0]);
      }

      this.setState({
        Date: event.value,
        placeholderTranslate: t('timepicker.placeholder'),
      });

      const { callBack } = this.props;
      const { Date } = this.state;

      if (callBack) {
        const filter = {
          timestamp_offset_start: 0,
          timestamp_offset_end: Math.floor(Date.getTime() / 1000),
          number_of_elements: 12,
          isScrolling: false,
          open: false,
          currentRecording: '',
        };
        callBack(filter);
      }
    }
  };

  render() {
    const { placeholderTranslate, Date } = this.state;
    return (
      <DateTimePickerComponent
        placeholder={placeholderTranslate}
        id="datetimepicker"
        strictMode="true"
        max={this.maxDate}
        onChange={this.handleChange}
        value={Date}
      />
    );
  }
}

TimePicker.propTypes = {
  callBack: PropTypes.func.isRequired,
};
export default TimePicker;
