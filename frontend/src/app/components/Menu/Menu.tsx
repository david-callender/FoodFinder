import { LikeButton } from "../LikeButton/LikeButton";

// components must be of shape FC
import type { FC } from "react";

// type for menu item
export type MenuItem = {
  meal: string;
  key?: number;
  location?: string;
};

type Props = {
  items: MenuItem[];
};

export const Menu: FC<Props> = ({ items }) => {
  // Description : Table for menu items. Iterates through items and generates a row.
  // Args:
  // items : string[]

  return (
    <div className="flex justify-center">
      <table className="self-center rounded-xl border border-white">
        <thead>
          <tr>
            <th>food</th>
            <th>location</th>
          </tr>
        </thead>

        <tbody>
          {items.map((item: MenuItem) => (
            <tr key={item.key}>
              <td className="p-5">{item.meal}</td>
              <td className="p-5">{item.location}</td>
              <td className="p-5">
                <LikeButton item={item} />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
