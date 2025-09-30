"use client";

import { useRouter } from "next/navigation";

import type { FC, FormEvent } from "react";

// for casting when fetching from api
type User = {
  access_token: string;
};

export const LoginBox: FC = () => {
  // Description: email/password text field and login button w/ redirect

  // use for redirects
  const router = useRouter();

  async function handleSubmit(
    event: FormEvent<HTMLFormElement>
  ): Promise<void> {
    // prevents refresh of page
    event.preventDefault();

    // pulling form data
    const formData = new FormData(event.currentTarget);
    const email = formData.get("email");
    const password = formData.get("password");

    // TODO: input validation

    // 9/27: endpoint currently assumes all users are valid.
    // At some point, when validation is implemented,
    // update this to a redirect (?) for registration
    // if login fails

    // request to login endpoint
    // refresh_token cookie is set here
    const response = await fetch(
      new URL("/login", process.env.NEXT_PUBLIC_BACKEND_URL),
      {
        method: "POST",
        credentials: "include", // need this for receive cookies w/ cors
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username: email, password: password }),
      }
    );

    if (response.ok) {
      // casting to known type w/ known fields
      const responseJson = (await response.json()) as User;

      // store access_token in local storage
      localStorage.setItem("access_token", responseJson.access_token);
      // redirect
      router.push("/");
    } else {
      console.log(await response.json());
    }
  }

  // final login box component
  return (
    <form onSubmit={handleSubmit}>
      <div className="m-20 flex flex-auto flex-col p-5">
        <input
          type="email"
          name="email"
          placeholder="Email"
          className="m-2 place-self-center rounded-lg border-4 bg-gray-200 p-0.5 text-black"
          required
        />
        <input
          type="password"
          name="password"
          placeholder="Password"
          className="m-2 place-self-center rounded-lg border-4 bg-gray-200 p-0.5 text-black"
          required
        />
        <button className="mx-auto w-40 bg-gray-200 text-black hover:bg-gray-300">
          login
        </button>
        <p className="m-2 place-self-center text-xs">
          Don&apos;t have an account? <a href="/signup">Sign Up Here</a>
        </p>
      </div>
    </form>
  );
};
