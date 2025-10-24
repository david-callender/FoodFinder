"use client";

import type { FC, FormEvent } from "react";

type Props = {
  setDiningHall: (diningHall: string) => void;
  setDate: (date: string | undefined) => void;
  setTime: (time: "breakfast" | "lunch" | "dinner" | "everyday") => void;
  handleSubmit: (e: FormEvent<HTMLFormElement>) => Promise<void>;
};

export const MealSearch: FC<Props> = ({
  setDiningHall,
  setDate,
  setTime,
  handleSubmit,
}) => {
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
                console.log(e.target.value);
                const date = e.target.value;
                if (date.length === 0) {
                  setDate(undefined);
                } else {
                  setDate(date);
                }
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
            <button className="m-2">Search</button>
          </div>
        </form>
      </div>
    </>
  );
};
