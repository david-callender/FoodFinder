import type { MenuItem } from "./components/Menu/Menu";

// sample data for display

const item1: MenuItem = {
  meal: "Apples",
  key: 5,
  location: "comstock",
};

const item2: MenuItem = {
  meal: "Bananas",
  key: 6,
  location: "pioneer",
};

export const MenuItems: MenuItem[] = [item1, item2];
