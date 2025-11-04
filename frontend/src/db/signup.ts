// see login.ts fo why is this commented out
// "use server";

import * as z from "zod";

const SCHEMA = z.object({
  accessToken: z.string(),
});

type SignupData = z.output<typeof SCHEMA>;

type SuccessResponse = { ok: true; data: SignupData };
type ErrorResponse = { ok: false; error: string };

type Response = SuccessResponse | ErrorResponse;

export const signup = async (
  email: string,
  password: string,
  displayName: string
): Promise<Response> => {
  // Purpose : Creating new user credentials given fields email, password
  // Args:
  // email : string - user's email
  // password: string - user's password
  // Returns:
  // {access_token: string} - user's access token

  // TODO [backend] : this will switch to signup when backend is fixed
  const registerURL = new URL("/signup", process.env.NEXT_PUBLIC_BACKEND_URL);

  const response = await fetch(registerURL, {
    method: "POST",
    credentials: "include", // need this for receive cookies w/ cors
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password, displayName }),
  });

  const json = (await response.json()) as unknown;

  return response.ok
    ? { ok: true, data: await SCHEMA.parseAsync(json) }
    : { ok: false, error: "Call to /signup failed: " + JSON.stringify(json) };
};
