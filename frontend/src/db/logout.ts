"use server";

export const logout = async (): Promise<void> => {
  // Purpose : Revoking current credentials for user session using "/logout" endpoint
  // Args:
  // None
  // Returns:
  // void

  const logoutURL = new URL("/logout", process.env.NEXT_PUBLIC_BACKEND_URL);
  // TODO [misc.] : Maybe revoke cookies?
  const response = await fetch(logoutURL, {
    method: "POST",
  });
  if (!response.ok) {
    const json = (await response.json()) as unknown;
    throw new Error("Call to /logout failed: " + JSON.stringify(json));
  }
};
