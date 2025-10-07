import type { MenuItem } from "../components/Menu/Menu";

// sample data for display

const item1: MenuItem = {
  id: 5,
  meal: "Apples",
  isPreferred: false,
};

const item2: MenuItem = {
  id: 6,
  meal: "Bananas",
  isPreferred: true,
};

export const MenuItems: MenuItem[] = [item1, item2];
