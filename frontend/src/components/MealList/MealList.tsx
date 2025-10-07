import { LikeButton } from "../LikeButton/LikeButton";

import type { MenuItem } from "../Menu/Menu";
import type { FC } from "react";

type Props = {
  items: MenuItem[];
};

export const MealList: FC<Props> = ({ items }) => {
  {
    return (
      <>
        {items.map((item: MenuItem) => (
          <tr key={item.id}>
            <td className="p-5">{item.meal}</td>
            <td className="p-5">
              <LikeButton item={item} />
            </td>
          </tr>
        ))}
      </>
    );
  }
};
