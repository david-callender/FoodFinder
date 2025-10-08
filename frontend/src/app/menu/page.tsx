import { Menu } from "../../components/Menu/Menu";
import { MenuItems } from "../seedData";

import type { FC } from "react";

export const Menu_Page: FC = () => {
  return (
    <>
      <div className="m-10 grid justify-center">
        <form>
          <div className="grid grid-cols-1">
            <input type="text" className="m-2 w-50 bg-white text-black"></input>
            <input type="date" className="m-2 w-50 bg-white text-black"></input>
            <input type="time" className="m-2 w-50 bg-white text-black"></input>
          </div>
        </form>
      </div>
      <Menu items={MenuItems} />
    </>
  );
};

export default Menu_Page;
