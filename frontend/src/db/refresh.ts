import { z } from "zod";

const SCHEMA = z.object({ accessToken: z.string() });

export const refresh = async (): Promise<string | undefined> => {
  const response = await fetch(
    new URL("/refresh", process.env.NEXT_PUBLIC_BACKEND_URL),
    {
      method: "POST",
      credentials: "include",
    }
  );

  const json = (await response.json()) as unknown;
  if (response.ok) {
    const { accessToken } = await SCHEMA.parseAsync(json);
    return accessToken;
  } else {
    alert(JSON.stringify(json));
    return undefined;
  }
};
