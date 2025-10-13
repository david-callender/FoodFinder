"use client";

import { useState } from "react";

import { MealSearch } from "@/components/MealSearch/MealSearch";

import { Menu } from "../../components/Menu/Menu";
import { MenuItems } from "../seedData";

import type { FC, FormEvent } from "react";

// for handling input from mealQuery form
type MealQuery = {
  meal: string;
  date: string;
  time: string;
};

// used for onSubmit behavior for MealQuery component. should match handleSubmit function
export type HandleMealQueryFunction = (
  event: FormEvent<HTMLFormElement>
) => void;

export const Menu_Page: FC = () => {
  // need to store this here so Menu Component can acces the future response object
  const [mealQuery, setMealQuery] = useState<MealQuery>();

  function handleSubmit(event: FormEvent<HTMLFormElement>): void {
    // Description : handling submit of search query from MealSearchComponent
    // TODO : make DB query to return list of items for user preferences, etc.
    // Args: MenuItem[]
    // items : string[]

    const formData = new FormData(event.currentTarget);

    // fetching form dat
    let meal = formData.get("meal");
    let date = formData.get("date");
    let time = formData.get("time");

    // Needed for validation against files (?). idk, it makes the linter happy
    // all fields are marked as required as well.
    meal = typeof meal === "string" ? meal : "MISSING_MEAL";
    date = typeof date === "string" ? date : "MISSING_DATE";
    time = typeof time === "string" ? time : "MISSING_TIME";

    // new meal query
    const newMealQuery: MealQuery = {
      meal: meal,
      date: date,
      time: time,
    };

    setMealQuery(newMealQuery);
    // logging meal query to make linter happy
    console.log(mealQuery);
  }

  return (
    <>
      <MealSearch handleSubmit={handleSubmit} />
      {/* TODO : update MenuItems from static -> response object from backend db
          MenuItems should be a state variable which is updated in handlSubmit */}
      <Menu items={MenuItems} />
    </>
  );
};

export default Menu_Page;
