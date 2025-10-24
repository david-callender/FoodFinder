"use client";

import Image from "next/image";
import { useState } from "react";

import { addFoodPreference } from "@/db/addFoodPreference";
import { removeFoodPreference } from "@/db/removeFoodPreference";

import filled_heart from "./filled.png";
import unfilled_heart from "./unfilled.png";

import type { SetPreferenceFunction } from "../Menu/Menu";
import type { MenuItem } from "@/db/getMenu";
import type { FC } from "react";

type Props = {
  item: MenuItem;
  handlePreferenceChange: SetPreferenceFunction;
};

export const LikeButton: FC<Props> = ({ item, handlePreferenceChange }) => {
  // Description : like button in menu table

  // setting state
  const [liked, setLiked] = useState(item.isPreferred);

  // updating liked value
  async function postLike(item: MenuItem): Promise<void> {
    const dbCall = item.isPreferred ? removeFoodPreference : addFoodPreference;

    await dbCall(item.meal);

    // automatically setting like button state to opposite of current state.
    // Should not need to call this function otherwise.
    setLiked(!item.isPreferred);
    // setting preference in parent Menu component on item.
    handlePreferenceChange(item);
  }

  // like button
  const like = (
    <button
      onClick={async () => {
        await postLike(item);
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
      onClick={async () => {
        await postLike(item);
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
