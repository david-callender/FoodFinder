// converts components to client components
"use client";

import Image from "next/image";
// for state management in components
import { useState } from "react";

// needed for image
// TODO: should be storing these images in db?
import filled_heart from "./filled.png";
import unfilled_heart from "./unfilled.png";

import type { MenuItem } from "../Menu/Menu";
// components must be of shape FC
import type { FC } from "react";

type Props = {
  item: MenuItem;
};

export const LikeButton: FC<Props> = ({ item }) => {
  // Description : like button in menu table

  // setting state
  const [liked, setLiked] = useState(item.isPreferred);

  // updating liked value
  function postLike(item: MenuItem, value: boolean): void {
    // TODO : make a POST request for the current user to like the food
    setLiked(value);
    console.log("liked " + item.meal);
  }

  // like button
  const like = (
    <button
      onClick={() => {
        postLike(item, false);
      }}
    >
      <Image
        src={filled_heart}
        alt="filled heart"
        className="w-15 transition duration-150 active:scale-90"
      />
    </button>
  );

  // not liked button
  const notLiked = (
    <button
      onClick={() => {
        postLike(item, true);
      }}
    >
      <Image
        src={unfilled_heart}
        alt="unfilled heart"
        className="w-15 transition duration-150 active:scale-90"
      />
    </button>
  );

  return liked ? like : notLiked;
};
