// for state management in components
import Image from "next/image";
import { useState } from "react";

// needed for image
import filled_heart from "./filled.png";
import "./LikeButton.css";
import unfilled_heart from "./unfilled.png";

// components must be of shape FC
import type { MenuProp } from "../Menu/Menu";
import type { FC } from "react";

export const LikeButton: FC<Pick<MenuProp, "MenuItem">> = ({ MenuItem }) => {
  // Description : like button in menu table

  // setting state
  const [liked, setLiked] = useState(false);

  // updating liked value
  function postLike(food: string, value: boolean): void {
    // TODO : make a POST request for the current user to like the food
    setLiked(value);
    console.log("liked" + food);
  }

  // like button
  const like = (
    <button
      onClick={() => {
        postLike(MenuItem, false);
      }}
    >
      <Image src={filled_heart} height={50} width={50} alt="filled heart" />
    </button>
  );

  // not liked button
  const notLiked = (
    <button
      onClick={() => {
        postLike(MenuItem, true);
      }}
    >
      <Image src={unfilled_heart} height={50} width={50} alt="unfilled heart" />
    </button>
  );

  return liked ? like : notLiked;
};
