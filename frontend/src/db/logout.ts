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
    throw new Error("UNABLE TO LOGOUT");
  }
};
