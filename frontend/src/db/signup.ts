// see login.ts fo why is this commented out
// "use server";

export const signup = async (
  email: string,
  password: string,
  displayName: string
): Promise<void> => {
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

  if (!response.ok) {
    const json = (await response.json()) as unknown;
    throw new Error("Call to /signup failed: " + JSON.stringify(json));
  }
};
