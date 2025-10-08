import { LikeButton } from "../LikeButton/LikeButton";

import type { MenuItem, setPreferenceFunction } from "../Menu/Menu";
import type { FC } from "react";

type Props = {
  items: MenuItem[];
  setPreference: setPreferenceFunction;
};

export const MealList: FC<Props> = ({ items, setPreference }) => {
  {
    return (
      <>
        {items.map((item: MenuItem) => (
          <div key={item.id} className="grid grid-cols-2 items-center">
            <div className="p-5">{item.meal}</div>
            <div className="p-5">
              <LikeButton item={item} setPreference={setPreference} />
            </div>
          </div>
        ))}
      </>
    );
  }
};
