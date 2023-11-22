import { useState } from "react";

import { Icon } from "~/res/icons/icon";

interface Props {
  rating: number;
  setRating: (rating: number) => void;
}

// https://dev.to/michaelburrows/create-a-custom-react-star-rating-component-5o6
export const Rating = (props: Props) => {
  const [hover, setHover] = useState(0);

  const onStarClick = (rating: number) => () => {
    props.setRating(rating);
  };

  const onMouseEnter = (rating: number) => () => {
    setHover(rating);
  };

  const onMouseLeave = () => {
    setHover(props.rating);
  };

  return (
    <>
      {[...Array(10)].map((_, index) => {
        index += 1;
        const icon = index <= (hover || props.rating) ? "star" : "star-outline";
        const fill =
          index <= (hover || props.rating) ? "fill-white" : "fill-slate-400";
        return (
          <button
            type="button"
            key={index}
            // className={index <= (hover || props.rating) ? 'on' : 'off'}
            onClick={onStarClick(index)}
            onMouseEnter={onMouseEnter(index)}
            onMouseLeave={onMouseLeave}
          >
            <Icon name={icon} size={18} fill={fill} />
          </button>
        );
      })}
    </>
  );
};
