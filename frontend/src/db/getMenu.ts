import { redirect } from "next/navigation";
import * as z from "zod";

import { err, ERROR_SCHEMA, ok } from "./error";
import { refresh } from "./refresh";

import type { Result } from "./error";

const SCHEMA = z.array(
  z.object({
    id: z.string(),
    meal: z.string(),
    isPreferred: z.boolean(),
  })
);

type MenuData = z.output<typeof SCHEMA>;

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
): Promise<Result<MenuData, string>> => {
  const accessToken = await refresh();

  if (accessToken === undefined) {
    redirect("/login");
  }

  const searchParams = new URLSearchParams({
    accessToken,
    day,
    mealtime,
    diningHall,
  });
  const response = await fetch(
    new URL("/getMenu?", process.env.NEXT_PUBLIC_BACKEND_URL).toString() +
      searchParams.toString()
  );
  const json = (await response.json()) as unknown;

  if (response.ok) {
    return ok(await SCHEMA.parseAsync(json));
  }

  const { detail } = await ERROR_SCHEMA.parseAsync(json);

  if (detail === "unauthenticated") {
    redirect("/login");
  }
  return err(detail);
};
