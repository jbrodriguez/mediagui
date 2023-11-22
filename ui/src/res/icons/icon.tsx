import React from "react";

interface IconProps {
  name: string;
  size?: number;
  fill?: string;
}

export const Icon: React.FC<IconProps> = ({ name, size = 24, fill = "" }) => {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" className={`${fill}`}>
      <path d={icons[name]} />
    </svg>
  );
};

const icons: { [key: string]: string } = {
  // Add your icon paths here
  binoculars:
    "M11,6H13V13H11V6M9,20A1,1 0 0,1 8,21H5A1,1 0 0,1 4,20V15L6,6H10V13A1,1 0 0,1 9,14V20M10,5H7V3H10V5M15,20V14A1,1 0 0,1 14,13V6H18L20,15V20A1,1 0 0,1 19,21H16A1,1 0 0,1 15,20M14,5V3H17V5H14Z",
  plus: "M19,13H13V19H11V13H5V11H11V5H13V11H19V13Z",
  star: "M12,17.27L18.18,21L16.54,13.97L22,9.24L14.81,8.62L12,2L9.19,8.62L2,9.24L7.45,13.97L5.82,21L12,17.27Z",
  "star-outline":
    "M12,15.39L8.24,17.66L9.23,13.38L5.91,10.5L10.29,10.13L12,6.09L13.71,10.13L18.09,10.5L14.77,13.38L15.76,17.66M22,9.24L14.81,8.63L12,2L9.19,8.63L2,9.24L7.45,13.97L5.82,21L12,17.27L18.18,21L16.54,13.97L22,9.24Z",
};
