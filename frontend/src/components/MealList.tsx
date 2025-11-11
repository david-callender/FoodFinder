import { LikeButton } from "./LikeButton";

import type { SetPreferenceFunction } from "./Menu";
import type { MenuItem } from "@/db/getMenu";
import type { FC } from "react";

type Props = {
  items: MenuItem[];
  handlePreferenceChange: SetPreferenceFunction;
};

export const MealList: FC<Props> = ({ items, handlePreferenceChange }) => {
  {
    return (
      <>
        {items.map((item: MenuItem) => (
          <div
            key={item.id}
            className="m-3 grid grid-cols-2 items-center rounded-2xl bg-black/40 p-4 shadow-sm transition-all duration-200 hover:bg-black/30 hover:shadow-md"
          >
            <div className="pl-2 text-lg font-medium text-white/90">
              {item.meal}
            </div>
            <div className="flex justify-end pr-2">
              <LikeButton
                item={item}
                handlePreferenceChange={handlePreferenceChange}
              />
            </div>
          </div>
        ))}
      </>
    );
  }
};
