import React from "react";

import DatePicker from "react-date-picker";
import "react-date-picker/dist/DatePicker.css";
import "react-calendar/dist/Calendar.css";

type ValuePiece = Date | null;
type Value = ValuePiece | [ValuePiece, ValuePiece];

export const Calendar = () => {
  const [value, setValue] = React.useState<Value>(new Date());

  return (
    <div className={`bg-white bg-clip-border text-black ml-3`}>
      <DatePicker
        onChange={setValue}
        value={value}
        clearIcon={null}
        format="y-MM-dd"
      />
    </div>
  );
};
