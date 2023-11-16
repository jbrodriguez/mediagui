import React from "react";

interface IconProps {
  name: string;
  size?: number;
  color?: string;
}

export const Icon: React.FC<IconProps> = ({
  name,
  size = 24,
  color = "currentColor",
}) => {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24">
      <path fill={color} d={icons[name]} />
    </svg>
  );
};

const icons: { [key: string]: string } = {
  // Add your icon paths here
  binoculars:
    "M11,6H13V13H11V6M9,20A1,1 0 0,1 8,21H5A1,1 0 0,1 4,20V15L6,6H10V13A1,1 0 0,1 9,14V20M10,5H7V3H10V5M15,20V14A1,1 0 0,1 14,13V6H18L20,15V20A1,1 0 0,1 19,21H16A1,1 0 0,1 15,20M14,5V3H17V5H14Z",
  plus: "M19,13H13V19H11V13H5V11H11V5H13V11H19V13Z",
};
