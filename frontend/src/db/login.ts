// if we use this, then the set-cookies headers are never accepted/
// I don't know what the benefits of having this here would be, so I'll just leave
// it commented out so that we have it for later.
// "use server";

import * as z from "zod";

const SCHEMA = z.object({
  displayName: z.string(),
  accessToken: z.string(),
});

type LoginData = z.output<typeof SCHEMA>;

// TODO! [misc.] : FIX THIS
// this is a bandaid solution until the /login endpoint is up and running
// zod will do the parsing of the response json, then the user will set their access token in local_storage
// for casting when fetching from api
export type User = {
  access_token: string;
};

export const login = async (
  email: string,
  password: string
): Promise<LoginData> => {
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
    body: JSON.stringify({ username: email, password: password }), // TODO : FIX This please! username should be email.
  });

  if (response.ok) {
    // TODO [misc.] : more fixing here (see export of type User, related)
    // casting to known type w/ known fields
    const responseJson = (await response.json()) as User;

    // TODO [backend] : stop returning dummy data
    return {
      displayName: "TEMPORARY_NAME",
      accessToken: responseJson.access_token,
    };
  } else {
    throw new Error("COULD NOT LOGIN");
  }

  // return username via zod
  const json = (await response.json()) as unknown;
  return await SCHEMA.parseAsync(json);
};
