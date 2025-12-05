import { handleError, ok } from "./error";

import type { Result } from "./error";

export const logout = async (): Promise<Result<undefined, string>> => {
  // Purpose : Revoking current credentials for user session using "/logout" endpoint
  // Args:
  // None
  // Returns:
  // void

  const logoutURL = new URL("/logout", process.env.NEXT_PUBLIC_BACKEND_URL);

  // TODO [misc.] : Maybe revoke cookies?
  const response = await fetch(logoutURL, {
    method: "POST",
    credentials: "include", // need this for receive cookies w/ cors
  });

  if (response.ok) {
    localStorage.removeItem("displayName");
    return ok(undefined);
  }

  const json = (await response.json()) as unknown;

  return await handleError(json);
};
