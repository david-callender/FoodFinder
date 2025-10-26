"use client";

import { useState } from "react";

import { getMenu } from "@/db/getMenu";

import { MealSearch, getCurrentDate } from "../MealSearch/MealSearch";
import { Menu } from "../Menu/Menu";

import type { MenuItem } from "@/db/getMenu";
import type { FC, FormEvent } from "react";

export const MenuManager: FC = () => {
  const [error, setError] = useState<string>("");
  // state for menu items
  const [menuItems, setMenuItems] = useState<MenuItem[]>([]);

  // state for query

  // diningHallID for comstock (default option in dropdown)
  const [diningHall, setDiningHall] = useState<string>(
    "62a90bbaa9f13a0e1cac2320"
  );
  const [date, setDate] = useState<string | undefined>(getCurrentDate());
  const [time, setTime] = useState<
    "breakfast" | "lunch" | "dinner" | "everyday"
    // set default to breakfast here because it's the first option in our drop down menu
    // and if we don't, and the user never selects a mealtime, then the field
    // is null/empty string in our query
  >("breakfast");

  async function handleSubmit(e: FormEvent<HTMLFormElement>): Promise<void> {
    // Description : handling submit of search query
    e.preventDefault();
    if (date === undefined) {
      setError("date was undefined");
      return;
    }

    setMenuItems(await getMenu(date, time, diningHall));
  }

  return (
    <>
      <div className="flex justify-center">
        <div className="flex w-1/2 flex-col">
          <MealSearch
            setDiningHall={setDiningHall}
            setDate={setDate}
            setTime={setTime}
            handleSubmit={handleSubmit}
          />
          <p>{error}</p>
          <div className="rounded-3xl bg-white/10 p-12">
            <Menu items={menuItems} />
          </div>
        </div>
      </div>
    </>
  );
};
