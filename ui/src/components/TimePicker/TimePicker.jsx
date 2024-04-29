import { DateTimePickerComponent } from '@syncfusion/ej2-react-calendars';
import React from 'react';
import './TimePicker.scss';
import { t } from 'i18next';

class TimePicker extends React.PureComponent {
  maxDate = new Date(new Date());
  
  constructor(props) {
    super(props);
    this.state = {
       Date: new Date(new Date()),
    };
    this.handleChange = this.handleChange.bind(this);
 }

  handleChange = (event) => {
    this.setState({
      Date: event.value
    })
  };

  render() {
    return (
      <DateTimePickerComponent
        placeholder={t('timepicker.placeholder')}
        id="datetimepicker"
        strictMode="true"
        max={this.maxDate}
        onChange={this.handleChange}
        value = {this.Date}
      />
    );
  }
}
export default TimePicker;
