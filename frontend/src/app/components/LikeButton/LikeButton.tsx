// for state management in components
import { useState } from "react";

// components must be of shape FC
import type { MenuProp } from "../Menu/Menu";
import type { FC } from "react";

export const LikeButton: FC<Pick<MenuProp, "MenuItem">> = ({
  MenuItem,
}: MenuProp) => {
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
      Liked
    </button>
  );

  // not liked button
  const notLiked = (
    <button
      onClick={() => {
        postLike(MenuItem, true);
      }}
    >
      Like
    </button>
  );

  return liked ? like : notLiked;
};
