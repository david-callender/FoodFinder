import { LikeButton } from "../LikeButton/LikeButton";

import type { SetPreferenceFunction } from "../Menu/Menu";
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
            className="grid grid-cols-2 items-center rounded-2xl bg-black/40 p-4 m-3 shadow-sm hover:shadow-md transition-all duration-200 hover:bg-black/30"
          >
            <div className="text-lg font-medium text-white/90 pl-2">
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
