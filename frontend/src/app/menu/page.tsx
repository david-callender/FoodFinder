import { Menu } from "../../components/Menu/Menu";
import { MenuItems } from "../seedData";

import type { FC } from "react";


export const Menu_Page: FC = () => {
  return <Menu items={MenuItems} />;
};

export default Menu_Page;
