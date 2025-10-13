import type { HandleMealQueryFunction } from "@/app/menu/page";
import type { FC } from "react";

type Props = {
  handleSubmit: HandleMealQueryFunction;
};

export const MealSearch: FC<Props> = ({ handleSubmit }) => {
  return (
    <>
      <div className="m-10 grid justify-center">
        <form onSubmit={handleSubmit}>
          <div className="grid grid-cols-1">
            <input
              type="text"
              name="meal"
              className="m-2 w-50 bg-white text-black"
              required
            ></input>
            <input
              type="date"
              name="date"
              className="m-2 w-50 bg-white text-black"
              required
            ></input>
            <input
              type="time"
              name="time"
              className="m-2 w-50 bg-white text-black"
              required
            ></input>
            <button>Search</button>
          </div>
        </form>
      </div>
    </>
  );
};
