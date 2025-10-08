import { Menu } from "../../components/Menu/Menu";
import { MenuItems } from "../seedData";

import type { FC } from "react";

export const Menu_Page: FC = () => {
  return (
    <>
      <form>
        <input type="text"></input>
      </form>
      <Menu items={MenuItems} />
    </>
  );
};

export default Menu_Page;
