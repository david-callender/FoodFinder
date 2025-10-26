"use client";

import type { FC, FormEvent } from "react";

type Props = {
  setDiningHall: (diningHall: string) => void;
  setDate: (date: string | undefined) => void;
  setTime: (time: "breakfast" | "lunch" | "dinner" | "everyday") => void;
  handleSubmit: (e: FormEvent<HTMLFormElement>) => Promise<void>;
};

export function getCurrentDate(): string {
  const today = new Date();
  const year = today.getFullYear();
  const month = String(today.getMonth() + 1).padStart(2, "0");
  const day = String(today.getDate()).padStart(2, "0");

  return `${year}-${month}-${day}`;
}

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
      <div className="m-10 flex justify-center">
        <form
          onSubmit={handleSubmit}
          className="w-full max-w-md rounded-2xl bg-white/10 p-6 shadow-lg backdrop-blur-md transition-all duration-300 hover:shadow-xl"
        >
          <h2 className="mb-4 text-center text-2xl font-semibold text-white">
            Search Menu
          </h2>

          <div className="space-y-4">
            <select
              name="diningHall"
              onChange={(e) => {
                setDiningHall(e.target.value);
              }}
              className="w-full rounded-xl border border-white/30 bg-white/80 px-3 py-2 text-black placeholder-gray-500 focus:border-blue-400 focus:ring-2 focus:ring-blue-400 focus:outline-none"
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
              value={getCurrentDate()}
              required
              onChange={(e) => {
                const date = e.target.value;
                setDate(date.length === 0 ? undefined : date);
              }}
              className="w-full rounded-xl border border-white/30 bg-white/80 px-3 py-2 text-black focus:border-blue-400 focus:ring-2 focus:ring-blue-400 focus:outline-none"
            />

            <select
              name="time"
              onChange={(e) => {
                changeTime(e.target.value);
              }}
              className="w-full rounded-xl border border-white/30 bg-white/80 px-3 py-2 text-black focus:border-blue-400 focus:ring-2 focus:ring-blue-400 focus:outline-none"
            >
              <option value="breakfast">Breakfast</option>
              <option value="lunch">Lunch</option>
              <option value="dinner">Dinner</option>
              <option value="everyday">Everyday</option>
            </select>

            <button
              type="submit"
              className="w-full rounded-xl bg-blue-500 px-4 py-2 font-semibold text-white shadow-md transition-all duration-200 hover:bg-blue-600 hover:shadow-lg"
            >
              Search
            </button>
          </div>
        </form>
      </div>
    </>
  );
};
