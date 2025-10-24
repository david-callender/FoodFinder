"use client";

import { useState } from "react";

import { getMenu } from "@/db/getMenu";

import { MealSearch } from "../MealSearch/MealSearch";
import { Menu } from "../Menu/Menu";

import type { MenuItem } from "@/db/getMenu";
import type { FC, FormEvent } from "react";

export const MenuManager: FC = () => {
  const [error, setError] = useState<string>("");
  // state for menu items
  const [menuItems, setMenuItems] = useState<MenuItem[]>([]);

  // state for query
  const [diningHall, setDiningHall] = useState<string | undefined>();
  const [date, setDate] = useState<string | undefined>();
  const [time, setTime] = useState<
    "breakfast" | "lunch" | "dinner" | "everyday"
    // set default to breakfast here because it's the first option in our drop down menu
    // and if we don't, and the user never selects a mealtime, then the field
    // is null/empty string in our query
  >("breakfast");

  async function handleSubmit(e: FormEvent<HTMLFormElement>): Promise<void> {
    // Description : handling submit of search query
    e.preventDefault();
    if (diningHall === undefined || date === undefined) {
      setError("something was undefined");
      return;
    }

    setMenuItems(await getMenu(date, time, diningHall));
  }

  return (
    <>
      <MealSearch
        setDiningHall={setDiningHall}
        setDate={setDate}
        setTime={setTime}
        handleSubmit={handleSubmit}
      />
      <p>{error}</p>
      <Menu items={menuItems} />
    </>
  );
};
