"use client";

import { MealList } from "../MealList/MealList";

// components must be of shape FC
import type { FC } from "react";

// type for menu item
export type MenuItem = {
  id: number;
  meal: string;
  isPreferred: boolean;
};

type Props = {
  items: MenuItem[];
};

export const Menu: FC<Props> = ({ items }) => {
  // Description : Table for menu items. Iterates through items and generates a row.
  // Args: MenuItem[]
  // items : string[]

  const preferred = items.filter((item: MenuItem) => item.isPreferred);

  const notPreferred = items.filter((item: MenuItem) => !item.isPreferred);

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
            <MealList items={preferred} />
          </tbody>

          <thead>
            <tr>
              <th>Not Preferred</th>
              <th></th>
            </tr>
          </thead>

          <tbody className="rounded-xl border border-white">
            <MealList items={notPreferred} />
          </tbody>
        </table>
      </div>
    </>
  );
};
