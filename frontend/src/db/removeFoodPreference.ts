// Purpose : updating preferred state of a certain meal to true using the "/addFoodPreference" endpoint
// Args:
// meal : string - literal string representing the meal
// Returns:

import { redirect } from "next/navigation";

import { ERROR_SCHEMA } from "./error";
import { refresh } from "./refresh";

// void - posting data to server
export const removeFoodPreference = async (meal: string): Promise<void> => {
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
    return;
  }

  const json = (await response.json()) as unknown;

  const { detail } = await ERROR_SCHEMA.parseAsync(json);

  if (detail === "unauthenticated") {
    redirect("/login");
  }

  throw new Error("Call to /removeFoodPreference failed: " + detail);
};
