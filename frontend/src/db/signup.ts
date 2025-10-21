// see login.ts fo why is this commented out
// "use server";

import * as z from "zod";

// bandaid import until we get the backend fully operational
// (returning both access_token AND displayName)
import type { User } from "./login";

const SCHEMA = z.object({
  displayName: z.string(),
  accessToken: z.string(),
});

type SignupData = z.output<typeof SCHEMA>;

export const signup = async (
  email: string,
  password: string
): Promise<SignupData> => {
  // Purpose : Creating new user credentials given fields email, password
  // Args:
  // email : string - user's email
  // password: string - user's password
  // Returns:
  // {access_token: string} - user's access token

  // TODO [backend] : this will switch to signup when backend is fixed
  const registerURL = new URL("/register", process.env.NEXT_PUBLIC_BACKEND_URL);

  const response = await fetch(registerURL, {
    method: "POST",
    credentials: "include", // need this for receive cookies w/ cors
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username: email, password: password }), // TODO : please fix this asap. username -> email on backend
  });

  if (response.ok) {
    const responseJson = (await response.json()) as User;

    // TODO [backend] : stop returning dummy data
    return {
      displayName: "TEMPORARY_NAME",
      accessToken: responseJson.access_token,
    };
  } else {
    throw new Error("COULD NOT SIGNUP");
  }

  // return username via zod
  const json = (await response.json()) as unknown;
  return await SCHEMA.parseAsync(json);
};
