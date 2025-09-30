"use client";

import { useRouter } from "next/navigation";

import type { FC, FormEvent } from "react";

export const SignUpBox: FC = () => {
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
    const phoneNumber = formData.get("phone-number");

    // TODO: accept phoneNumber at /register endpoint
    console.log(phoneNumber);
    // TODO: input validation

    // request to login endpoint
    // refresh_token cookie is set here
    const response = await fetch(
      new URL("/register", process.env.NEXT_PUBLIC_BACKEND_URL),
      {
        method: "POST",
        credentials: "include", // need this for receive cookies w/ cors
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username: email, password: password }),
      }
    );

    if (response.ok) {
      router.push("/");
    } else {
      console.log(response.status);
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <div className="flex flex-col">
        <input
          type="email"
          className="m-2 place-self-center rounded-lg border-4 bg-gray-200 p-0.5 text-black"
          name="email"
          placeholder="Email"
          required
        />
        <input
          type="password"
          className="m-2 place-self-center rounded-lg border-4 bg-gray-200 p-0.5 text-black"
          name="password"
          placeholder="Password"
          required
        />
        <input
          type="tel"
          className="m-2 place-self-center rounded-lg border-4 bg-gray-200 p-0.5 text-black"
          name="phone-number"
          placeholder="XXX-XXX-XXXX"
          pattern="[0-9]{3}-[0-9]{3}-[0-9]{4}"
        />
        <button className="mx-auto w-40 bg-gray-200 text-black hover:bg-gray-300">
          Sign Up!
        </button>
        <p className="m-2 place-self-center text-xs">
          Already have an account? <a href="/login">Login Here</a>
        </p>
      </div>
    </form>
  );
};
