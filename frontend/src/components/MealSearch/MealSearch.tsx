"use client";

import { useState } from "react";

import { getMenu } from "@/db/getMenu";

import type { FC } from "react";

export const MealSearch: FC = () => {
  const [diningHall, setDiningHall] = useState<string>();
  const [date, setDate] = useState<Date>();
  const [time, setTime] = useState<
    "breakfast" | "lunch" | "dinner" | "everyday"
  >();
  const [error, setError] = useState<string>("");

  async function handleSubmit(): Promise<void> {
    // Description : handling submit of search query

    if (diningHall === undefined || date === undefined || time === undefined) {
      setError("something was undefined");
      return;
    }

    // TODO : make DB query to return list of items for user preferences, etc.

    await getMenu(date, time, diningHall);
  }

  function changeTime(timeString: string): void {
    switch (timeString) {
      case "breakfast": {
        setTime("breakfast");

        break;
      }
      case "lunch": {
        setTime("lunch");

        break;
      }
      case "dinner": {
        setTime("dinner");

        break;
      }
      default: {
        setTime("everyday");
      }
    }
  }

  return (
    <>
      <div className="m-10 grid justify-center">
        <form onSubmit={handleSubmit}>
          <div className="grid grid-cols-1">
            <input
              type="text"
              name="diningHall"
              onChange={(e) => {
                setDiningHall(e.target.value);
              }}
              className="m-2 w-50 rounded-xl border-4 bg-white p-0.5 text-black"
              required
            ></input>
            <input
              type="date"
              name="date"
              onChange={(e) => {
                setDate(new Date(e.target.value));
              }}
              className="m-2 w-50 rounded-xl border-4 bg-white p-0.5 text-black"
              required
            ></input>
            <select
              name="time"
              className="m-2 rounded-xl border-5 border-white bg-white p-0.5 text-black"
              onChange={(e) => {
                changeTime(e.target.value);
              }}
            >
              <option value="breakfast">Breakfast</option>
              <option value="lunch">Lunch</option>
              <option value="dinner">Dinner</option>
              <option value="everyday">Everyday</option>
            </select>
            <button>Search</button>
          </div>
        </form>
        <p>{error}</p>
      </div>
    </>
  );
};
