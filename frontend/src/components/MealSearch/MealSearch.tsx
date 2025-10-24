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
            <select
              name="diningHall"
              onChange={(e) => {
     
                setDiningHall(e.target.value);
              }}
              className="m-2 rounded-xl border-5 border-white bg-white p-0.5 text-black"
            >
              <option value="62a90bbaa9f13a0e1cac2320">Comstock</option>
              <option value="6262b663b63f1e1517b6e433">Pioneer</option>
              <option value="627bbf3bb63f1e0fb3c1691a">17th Ave</option>
              <option value="627bbf2cb63f1e10059b45a4">Sanford</option>
              <option value="627bbeb6b63f1e0fa1c9fe7b">Middlebrook</option>
              <option value="62b21c96a9f13a0ac1472ef1">Bailey</option>
            </select>
            <input
              type="date"
              name="date"
              onChange={(e) => {
           
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
