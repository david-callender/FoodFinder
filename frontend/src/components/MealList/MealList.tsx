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
          <tr key={item.id}>
            <td className="p-5">{item.meal}</td>
            <td className="p-5">
              <LikeButton item={item} setPreference={setPreference} />
            </td>
          </tr>
        ))}
      </>
    );
  }
};
