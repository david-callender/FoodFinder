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

export const Menu: FC<Props> = ({ items: items_in }) => {
  const [items, setItems] = useState(items_in);

  useEffect(() => {
    setItems(items_in);
  }, [items_in]);

  // filter preferred from items
  const preferred = items.filter((item: MenuItem) => item.isPreferred);

  // everything else
  const notPreferred = items.filter((item: MenuItem) => !item.isPreferred);

  const handlePreferenceChange = (item: MenuItem): void => {
    const index = items.findIndex((i) => i.id === item.id);
    if (index === -1) {
      return;
    }

    setItems([
      ...items.slice(0, index),
      { ...items[index], isPreferred: !items[index].isPreferred },
      ...items.slice(index + 1),
    ]);
  };

  return preferred.length === 0 && notPreferred.length === 0 ? (
    <></>
  ) : (
    <>
      <div className="rounded-3xl bg-white/10 p-12">
        <div className="grid justify-center">
          <div className="min-h-10 rounded-xl border border-white">
            {/* To understand what is happening here, I would highly recommend reading this:
              https://stackoverflow.com/questions/76958201/how-to-pass-props-to-child-component-in-next-js-13
              essentially passing state function down two child components so that it can be used by LikeButton */}
            <MealList
              items={preferred}
              handlePreferenceChange={handlePreferenceChange}
            />
          </div>

          <div className="min-h-10 rounded-xl border border-white">
            <MealList
              items={notPreferred}
              handlePreferenceChange={handlePreferenceChange}
            />
          </div>
        </div>
      </div>
    </>
  );
};
