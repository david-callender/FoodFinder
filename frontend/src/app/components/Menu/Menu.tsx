import { LikeButton } from "../LikeButton/LikeButton";

// css for menu
import "./menu.css";

// components must be of shape FC
import type { FC } from "react";

// type for menu item
export type MenuProp = {
  MenuItem: string;
  key?: number;
  location?: string;
};

// type for collection of menu items
// need so we can pass into a component
export type MenuProps = {
  items: MenuProp[];
};

export const Menu: FC<MenuProps> = ({ items }: MenuProps) => {
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
          {items.map((item: MenuProp) => (
            <tr key={item.key}>
              <td className="menu-item">{item.MenuItem}</td>
              <td className="menu-item">{item.location}</td>
              <td className="menu-item">
                <LikeButton MenuItem={item.MenuItem} />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
