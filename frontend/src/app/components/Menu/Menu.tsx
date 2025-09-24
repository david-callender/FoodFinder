import { LikeButton } from "../LikeButton/LikeButton";

// css for menu
import "./menu.css";

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
    <div className="menu">
      <table className="menu">
        <thead>
          <tr>
            <th>food</th>
            <th>location</th>
          </tr>
        </thead>

        <tbody>
          {items.map((item: MenuItem) => (
            <tr key={item.key}>
              <td className="menu-item">{item.meal}</td>
              <td className="menu-item">{item.location}</td>
              <td className="menu-item">
                <LikeButton item={item} />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
