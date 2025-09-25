import type { MenuItem } from "../components/Menu/Menu";

// sample data for display

const item1: MenuItem = {
  key: 5,
  meal: "Apples",
  location: "comstock",
};

const item2: MenuItem = {
  key: 6,
  meal: "Bananas",
  location: "pioneer",
};

export const MenuItems: MenuItem[] = [item1, item2];
