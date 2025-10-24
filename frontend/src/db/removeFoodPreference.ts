"use server";

export const removeFoodPreference = async (
  meal: string
  // could make this boolean to confirm on the user end?
): Promise<void> => {
  // Purpose : updating preferred state of a certain meal to true using the "/addFoodPreference" endpoint
  // Args:
  // meal : string - literal string representing the meal
  // Returns:
  // void - posting data to server

  const foodPreferenceURL = new URL(
    "/removeFoodPreference",
    process.env.NEXT_PUBLIC_BACKEND_URL
  );

  const response = await fetch(foodPreferenceURL, {
    method: "POST",
    body: JSON.stringify({ meal }),
  });

  if (!response.ok) {
    const json = (await response.json()) as unknown;
    throw new Error(
      "Call to /removeFoodPreference failed: " + JSON.stringify(json)
    );
  }
};
