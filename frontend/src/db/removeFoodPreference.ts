"use server";

export const removeFoodPreference = async (
  meal: string
  // could make this boolean to confirm on the user end?
): Promise<void> => {
  // Purpose : updating preferred state of a certain meal to true using the "/add_food_preference" endpoint
  // Args:
  // meal : string - literal string representing the meal
  // Returns:
  // void - posting data to server

  const foodPreferenceURL = new URL(
    "/remove_food_preference",
    process.env.NEXT_PUBLIC_BACKEND_URL
  );

  const response = await fetch(foodPreferenceURL, {
    method: "POST",
    body: JSON.stringify({ value: meal }),
  });

  if (!response.ok) {
    throw new Error("Invalid Request");
  }
};
