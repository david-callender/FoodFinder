import * as z from "zod";

import { handleError, ok } from "./error";

import type { Result } from "./error";

const SCHEMA = z.object({
  displayName: z.string(),
});

export type LoginData = z.output<typeof SCHEMA>;

export const login = async (
  email: string,
  password: string
): Promise<Result<LoginData, string>> => {
  // Purpose : Login into/check credentials of a user
  // Args:
  // email : string - users email
  // password : string - users password
  // Returns:
  // {displayName: string, accessToken: string} - username and access token for the current session

  const loginURL = new URL("/login", process.env.NEXT_PUBLIC_BACKEND_URL);
  const response = await fetch(loginURL, {
    method: "POST",
    credentials: "include", // need this for receive cookies w/ cors
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });

  const json = (await response.json()) as unknown;

  if (response.ok) {
    return ok(await SCHEMA.parseAsync(json));
  }

  return await handleError(json);
};
