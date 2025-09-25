// converts components to client components
// necessary for client interactivity
// specifically the like button, in this case
// TODO: figure out if this is necessary for the scope of the entire page
"use client";

import { Menu } from "../components/Menu/Menu";

import { MenuItems } from "./seedData";

// components must be of shape FC
import type { FC } from "react";

const Home: FC = () => {
  // TODO : make api calls here to construct menu list
  // this is only temporary data

  return <Menu items={MenuItems} />;
};

export default Home;
