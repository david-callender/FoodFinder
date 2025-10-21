"use server";

export const addFoodPreference = async (meal: string): Promise<void> => {
  // Purpose : POSTing meal, updating preferred state using the "/add_food_preference" endpoint
  // Args:
  // meal : string - literal string representing the meal
  // Returns
  // void - posting data to server

  const foodPreferenceURL = new URL(
    "/add_food_preference",
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
