// Purpose : updating preferred state of a certain meal to true using the "/addFoodPreference" endpoint
// Args:
// meal : string - literal string representing the meal
// Returns:

import { redirect } from "next/navigation";

import { handleError, ok } from "./error";
import { refresh } from "./refresh";

import type { Result } from "./error";

// void - posting data to server
export const removeFoodPreference = async (
  meal: string
): Promise<Result<undefined, string>> => {
  const accessToken = await refresh();

  if (accessToken === undefined) {
    redirect("/login");
  }

  const foodPreferenceURL = new URL(
    "/removeFoodPreference",
    process.env.NEXT_PUBLIC_BACKEND_URL
  );

  const response = await fetch(foodPreferenceURL, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify({ accessToken, meal }),
  });

  if (response.ok) {
    return ok(undefined);
  }

  const json = (await response.json()) as unknown;

  return await handleError(json);
};
