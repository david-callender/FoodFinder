import type { MenuItem } from "@/db/getMenu";

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

export const seedMenuItems: MenuItem[] = [item1, item2];
