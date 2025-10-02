import { LinkBar } from "../components/LinkBar/LinkBar";
import { Menu } from "../components/Menu/Menu";

import { MenuItems } from "./seedData";

// components must be of shape FC
import type { FC } from "react";

const Home: FC = () => {
  // TODO : make api calls here to construct menu list
  // this is only temporary data
  return (
    <>
      <LinkBar />
      <Menu items={MenuItems} />
    </>
  );
};

export default Home;
