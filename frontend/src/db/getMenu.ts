"use server";

import * as z from "zod";

const SCHEMA = z.array(
  z.object({
    id: z.string(),
    meal: z.string(),
    isPreferred: z.boolean(),
  })
);

type MenuData = z.output<typeof SCHEMA>;

type SuccessResponse = { ok: true; data: MenuData };
type ErrorResponse = { ok: false; error: string };

type Response = SuccessResponse | ErrorResponse;

export type MenuItem = MenuData[number];

// Purpose : retrieving data for meals given a set of parameters
// Args:
// date : Date - day for meal data. Must be in YYYY-MM-DD format
// mealtime : "breakfast" | "lunch" | "dinner" | "everyday" - meal string representing the time of day for the meal
// diningHall : string - which dining hall to query menu for
// Returns:
// {meal: string, isPreferred: bool, id: string}[] - list of meals that matched the given search criteria
export const getMenu = async (
  day: string,
  mealtime: "breakfast" | "lunch" | "dinner" | "everyday",
  diningHall: string
): Promise<Response> => {
  const searchParams = new URLSearchParams({
    day,
    mealtime,
    diningHall,
  });
  const response = await fetch(
    new URL("/getMenu?", process.env.NEXT_PUBLIC_BACKEND_URL).toString() +
      searchParams.toString()
  );
  const json = (await response.json()) as unknown;

  return response.ok
    ? { ok: true, data: await SCHEMA.parseAsync(json) }
    : { ok: false, error: "call to /getMenu failed: " + JSON.stringify(json) };
};
