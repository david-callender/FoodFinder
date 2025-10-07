"use client";

import { useState } from "react";

import { MealList } from "../MealList/MealList";

// components must be of shape FC
import type { FC } from "react";

// type for menu item
export type MenuItem = {
  id: number;
  meal: string;
  isPreferred: boolean;
};

// used for typing in child components see: MealList/LikeButton component
export type setPreferenceFunction = (item: MenuItem) => void;

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

  function handlePreferenceChange(item: MenuItem): void {
    // TODO : make a POST request for the current user to like the food
    // what the boolean is changing to
    const isPreferred = !item.isPreferred;

    // updating to preferred food
    if (isPreferred) {
      // remove from notPreferred
      setNotPreferred(
        notPreferred.filter((tempItem: MenuItem) => tempItem !== item)
      );
      // updating preference on the item
      item.isPreferred = isPreferred;
      // add to preferred
      setPreferred([...preferred, item]);
      // changing to not preferred
    } else {
      // remove from preferred
      setPreferred(preferred.filter((tempItem: MenuItem) => tempItem !== item));
      // updating preference on the item
      item.isPreferred = isPreferred;
      // add to notPreferred
      setNotPreferred([...notPreferred, item]);
    }
  }

  return (
    <>
      <div className="flex justify-center">
        <table className="self-center rounded-xl">
          <thead>
            <tr>
              <th>Preferred</th>
            </tr>
          </thead>

          <tbody className="rounded-xl border border-white">
            {/* To understand what is happening here, I would highly recommend reading this:
            https://stackoverflow.com/questions/76958201/how-to-pass-props-to-child-component-in-next-js-13
            essentially passing state function down two child components so that it can be used by LikeButton */}
            <MealList
              items={preferred}
              setPreference={handlePreferenceChange}
            />
          </tbody>

          <thead>
            <tr>
              <th>Not Preferred</th>
              <th></th>
            </tr>
          </thead>

          <tbody className="rounded-xl border border-white">
            <MealList
              items={notPreferred}
              setPreference={handlePreferenceChange}
            />
          </tbody>
        </table>
      </div>
    </>
  );
};
