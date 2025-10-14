import { LikeButton } from "../LikeButton/LikeButton";

import type { SetPreferenceFunction } from "../Menu/Menu";
import type { MenuItem } from "@/db/getMenu";
import type { FC } from "react";

type Props = {
  items: MenuItem[];
  setPreference: SetPreferenceFunction;
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
