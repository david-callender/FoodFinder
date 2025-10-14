"use server";

import * as z from "zod";

const SCHEMA = z.array(
  z.object({
    id: z.number(),
    meal: z.string(),
    isPreferred: z.boolean(),
  })
);

type MenuData = z.output<typeof SCHEMA>;

export type MenuItem = MenuData[number];

export const getMenu = async (
  day: Date,
  mealtime: "breakfast" | "lunch" | "dinner" | "everyday",
  diningHall: string
): Promise<MenuData> => {
  return [{ id: 20, isPreferred: true, meal: "bana" }];
  const searchParams = new URLSearchParams({
    day: day.toISOString(),
    mealtime: mealtime,
    dining_hall: diningHall,
  });
  const response = await fetch(
    new URL("/get_menu?", process.env.NEXT_PUBLIC_BACKEND_URL).toString() +
      searchParams.toString()
  );
  const json = (await response.json()) as unknown;
  return await SCHEMA.parseAsync(json);
};
