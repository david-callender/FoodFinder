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

    localStorage.setItem("access_token", response.accessToken);
    // redirect
    router.push("/");
  }

  // final login box component
  return (
    <form onSubmit={handleSubmit}>
      <div className="flex flex-col p-5">
        <input
          type="email"
          onChange={(e) => {
            setEmail(e.target.value);
          }}
          name="email"
          placeholder="Email"
          className="m-3 place-self-center rounded-lg border-4 border-gray-200 bg-gray-200 p-0.5 text-black"
          required
        />
        <input
          type="password"
          onChange={(e) => {
            setPassword(e.target.value);
          }}
          name="password"
          placeholder="Password"
          className="m-3 place-self-center rounded-lg border-4 border-gray-200 bg-gray-200 p-0.5 text-black"
          required
        />
        <button className="mx-auto w-40 bg-gray-200 text-black hover:bg-gray-300">
          login
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
