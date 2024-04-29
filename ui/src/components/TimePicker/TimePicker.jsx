import React from 'react';
import { DateTimePickerComponent } from '@syncfusion/ej2-react-calendars';
import './TimePicker.scss';
import { t } from 'i18next';

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
    this.setState({
      Date: event.value,
      placeholderTranslate: t('timepicker.placeholder'),
    });
    const { callBack } = this.props;

    if (callBack) {
      const filter = {
        timestamp_offset_start: 0,
        timestamp_offset_end: Math.floor(this.state.Date.getTime() / 1000),
        number_of_elements: 12,
        isScrolling: false,
        open: false,
        currentRecording: '',
      };
      callBack(filter);
    }
  };

  render() {
    return (
      <DateTimePickerComponent
        placeholder={this.state.placeholderTranslate}
        id="datetimepicker"
        strictMode="true"
        max={this.maxDate}
        onChange={this.handleChange}
        value={this.state.Date}
      />
    );
  }
}
export default TimePicker;
