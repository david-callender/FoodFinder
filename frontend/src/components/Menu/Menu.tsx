"use client";

import { useEffect, useState } from "react";

import { MealList } from "../MealList/MealList";

import type { MenuItem } from "@/db/getMenu";
// components must be of shape FC
import type { FC } from "react";

// used for typing in child components see: MealList/LikeButton component
export type SetPreferenceFunction = (item: MenuItem) => void;

type Props = {
  items: MenuItem[];
};

export const Menu: FC<Props> = ({ items }) => {
  // Description : Table for menu items. Iterates through items and generates a row.
  // Args: MenuItem[]
  // items : string[]

  // filter preferred from items
  const basePreferred = items.filter((item: MenuItem) => item.isPreferred);

  // everything else
  const baseNotPreferred = items.filter((item: MenuItem) => !item.isPreferred);

  // keep two seperate states for preferred/not preferred
  const [preferred, setPreferred] = useState(basePreferred);
  const [notPreferred, setNotPreferred] = useState(baseNotPreferred);

  useEffect(() => {
    // filter preferred from items
    const basePreferred = items.filter((item: MenuItem) => item.isPreferred);

    // everything else
    const baseNotPreferred = items.filter(
      (item: MenuItem) => !item.isPreferred
    );

    setPreferred(basePreferred);
    setNotPreferred(baseNotPreferred);
  }, [items]);

  function handlePreferenceChange(item: MenuItem): void {
    //  TODO [db] : make a POST request for the current user to like the food
    // flipping preference status
    const isPreferred = !item.isPreferred;

    // moving to preferred food
    if (isPreferred) {
      // remove from notPreferred
      setNotPreferred(
        // TODO [misc.] : is this the comparison to be made? Not quite sure if object comparison like this is bulletproof
        notPreferred.filter((tempItem: MenuItem) => tempItem !== item)
      );
      // updating preference on the item
      item.isPreferred = isPreferred;
      setPreferred([...preferred, item]);
      // TODO [backend] : leave commented out until fully implemented in backend.
      // Also see removeFoodPreference in the other branch of this if statement
      //addFoodPreference(item.meal);

      // moving to not preferred
    } else {
      // remove from preferred
      setPreferred(preferred.filter((tempItem: MenuItem) => tempItem !== item));
      // updating preference on the item
      item.isPreferred = isPreferred;
      setNotPreferred([...notPreferred, item]);
      //removeFoodPreference(item.meal);
    }
  }

  return (
    <>
      <div className="grid justify-center">
        <div className="min-h-10 rounded-xl border border-white">
          {/* To understand what is happening here, I would highly recommend reading this:
            https://stackoverflow.com/questions/76958201/how-to-pass-props-to-child-component-in-next-js-13
            essentially passing state function down two child components so that it can be used by LikeButton */}
          <MealList items={preferred} setPreference={handlePreferenceChange} />
        </div>

        <div className="min-h-10 rounded-xl border border-white">
          <MealList
            items={notPreferred}
            setPreference={handlePreferenceChange}
          />
        </div>
      </div>
    </>
  );
};
