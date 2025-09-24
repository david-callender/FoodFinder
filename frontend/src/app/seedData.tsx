import type { MenuProp, MenuProps } from "./page";

// sample data for display

const item1: MenuProp = {
  MenuItem: "Apples",
  key: 5,
  location: "comstock",
};

const item2: MenuProp = {
  MenuItem: "Bananas",
  key: 6,
  location: "pioneer",
};

export const MenuItems: MenuProps = {
  items: [item1, item2],
};
