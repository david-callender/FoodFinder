// if we use this, then the set-cookies headers are never accepted/
// I don't know what the benefits of having this here would be, so I'll just leave
// it commented out so that we have it for later.
// "use server";

import * as z from "zod";

const SCHEMA = z.object({
  displayName: z.string(),
  accessToken: z.string(),
});

export type LoginData = z.output<typeof SCHEMA>;

type SuccessResponse = { ok: true; data: LoginData };
type ErrorResponse = { ok: false; error: string };

type Response = SuccessResponse | ErrorResponse;

export const login = async (
  email: string,
  password: string
): Promise<Response> => {
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

  return response.ok
    ? { ok: true, data: await SCHEMA.parseAsync(json) }
    : { ok: false, error: "call to /getMenu failed: " + JSON.stringify(json) };
};
