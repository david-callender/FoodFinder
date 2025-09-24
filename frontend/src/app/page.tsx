// converts components to client components
// necessary for client interactivity
// specifically the like button, in this case
// TODO: figure out if this is necessary for the scope of the entire page
"use client";

// for state management in components
import { useState } from "react";

import { MenuItems } from "./seedData";

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

const LikeButton: FC<Pick<MenuProp, "MenuItem">> = ({ MenuItem }: MenuProp) => {

  // Description : like button in menu table

  // setting state
  const [liked, setLiked] = useState(false);

  // updating liked value
  function postLike(food: string, value: boolean): void {
    // TODO : make a POST request for the current user to like the food
    setLiked(value);
    console.log("liked" + food);
  }

  // like button
  const like = (
    <button
      onClick={() => {
        postLike(MenuItem, false);
      }}
    >
      Liked
    </button>
  );

  // not liked button
  const notLiked = (
    <button
      onClick={() => {
        postLike(MenuItem, true);
      }}
    >
      Like
    </button>
  );

  return liked ? like : notLiked;
};

const Menu: FC<MenuProps> = ({ items }: MenuProps) => {

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

const Home: FC = () => {
  // make api calls here to construct menu list
  // this is only temporary data

  return <Menu items={MenuItems.items} />;
};

export default Home;
