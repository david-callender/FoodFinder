// converts components to client components
"use client";

import Image from "next/image";
// for state management in components
import { useState } from "react";

// needed for image
// TODO: should be storing these images in db?
import filled_heart from "./filled.png";
import unfilled_heart from "./unfilled.png";

import type { MenuItem, SetPreferenceFunction } from "../Menu/Menu";
// components must be of shape FC
import type { FC } from "react";

type Props = {
  item: MenuItem;
  setPreference: SetPreferenceFunction;
};

export const LikeButton: FC<Props> = ({ item, setPreference }) => {
  // Description : like button in menu table

  // setting state
  const [liked, setLiked] = useState(item.isPreferred);

  // updating liked value
  function postLike(item: MenuItem): void {
    // automatically setting like button state to opposite of current state.
    // Should not need to call this function otherwise.
    setLiked(!item.isPreferred);
    // setting preference in parent Menu component on item.
    setPreference(item);
  }

  // like button
  const like = (
    <button
      onClick={() => {
        postLike(item);
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
        postLike(item);
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
