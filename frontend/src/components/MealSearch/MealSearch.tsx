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
              className="m-2 w-50 bg-white text-black rounded-xl border-4 p-0.5"
              required
            ></input>
            <input
              type="date"
              name="date"
              className="m-2 w-50 bg-white text-black rounded-xl border-4 p-0.5"
              required
            ></input>
            <select name="time" className="bg-white m-2 text-black rounded-xl border-5 border-white p-0.5">
                <option value="0">Breakfast</option>
                <option value="1">Lunch</option>
                <option value="2">Dinner</option>
                <option value="3">Everyday</option>
            </select>
            <button>Search</button>
          </div>
        </form>
      </div>
    </>
  );
};
