"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";

import { login } from "@/db/login";

import type { FC, FormEvent } from "react";

export const LoginBox: FC = () => {
  // Description: email/password text field and login button w/ redirect

  // use for redirects
  const router = useRouter();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  async function handleSubmit(
    event: FormEvent<HTMLFormElement>
  ): Promise<void> {
    // prevents refresh of page
    event.preventDefault();

    // 9/27: endpoint currently assumes all users are valid.
    // At some point, when validation is implemented,
    // update this to a redirect (?) for registration
    // if login fails

    // request to login endpoint
    // refresh_token cookie is set here
    const response = await login(email, password);

    localStorage.setItem("displayName", response.displayName);

    // redirect
    router.push("/");
  }

  // final login box component
  return (
    <form onSubmit={handleSubmit}>
      <div className="flex flex-col p-5">
        <input
          type="email"
          value={email}
          onChange={(e) => {
            setEmail(e.target.value);
          }}
          name="email"
          placeholder="Email"
          className="m-3 w-60 place-self-center rounded-lg bg-gray-800 p-2"
          required
        />
        <input
          type="password"
          value={password}
          onChange={(e) => {
            setPassword(e.target.value);
          }}
          name="password"
          placeholder="Password"
          className="m-3 w-60 place-self-center rounded-lg bg-gray-800 p-2"
          required
        />
        <button className="mx-auto rounded-xl bg-red-900 px-4 py-2 font-semibold shadow transition hover:cursor-pointer hover:bg-red-700">
          Login
        </button>
        <p className="m-2 place-self-center text-xs">
          Don&apos;t have an account?{" "}
          <a href="/signup" className="text-blue-400">
            Sign Up Here
          </a>
        </p>
      </div>
    </form>
  );
};
