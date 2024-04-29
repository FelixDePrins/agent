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
  }

  handleChange = (event) => {
    this.setState({
      Date: event.value,
      placeholderTranslate: t('timepicker.placeholder'),
    });
    const { callBack } = this.props;
    if (callBack) {
      callBack();
    }
  };

  render() {
    return (
      <DateTimePickerComponent
        placeholder={this.state.placeholderTranslate}
        id="datetimepicker"
        strictMode={true}
        max={this.maxDate}
        onChange={this.handleChange}
        value={this.state.Date}
      />
    );
  }
}
export default TimePicker;
