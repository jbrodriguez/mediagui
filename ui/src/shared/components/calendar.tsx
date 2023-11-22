import React from "react";

import DatePicker from "react-date-picker";
import "react-date-picker/dist/DatePicker.css";
import "react-calendar/dist/Calendar.css";
import formatRFC3339 from "date-fns/formatRFC3339";

type ValuePiece = Date | null;
type Value = ValuePiece | [ValuePiece, ValuePiece];

interface Props {
  onChange: (value: string) => void;
}

export const Calendar: React.FC<Props> = (props) => {
  const [value, setValue] = React.useState<Value>(new Date());

  const onChange = (value: Value) => {
    setValue(value);
    const watched = formatRFC3339(value as Date);
    props.onChange(watched);
  };

  return (
    <div className={`bg-white bg-clip-border text-black ml-3`}>
      <DatePicker
        onChange={onChange}
        value={value}
        clearIcon={null}
        format="y-MM-dd"
      />
    </div>
  );
};
