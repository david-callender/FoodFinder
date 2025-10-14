"use client";

import { useRouter } from "next/navigation";

import { PhoneNumberInput } from "../PhoneNumberInput/PhoneNumberInput";

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

    // hack-y solution for now. Because the valid string  in the input component contains special characters for masking,
    // we end up pulling thos special characters with formData.get() have to replace them.

    // in an ideal world, we could "hijack" the get() function to return some value that isn't the value inside
    // of the input component and instead return the phone number string directly.

    let phoneNumber = formData.get("phone-number");

    // if phoneNumber is entered
    if (phoneNumber !== "") {
      // clean phone number of masking characters used in <PhoneNumberInput />
      phoneNumber = phoneNumber as string;
      phoneNumber = phoneNumber.replaceAll(/[^0-9]/g, "");
      // handle phone number db stuff here
    }

    console.log(phoneNumber);

    // TODO: DB accept phoneNumber at /register endpoint

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

    // TODO : Handling a "user exists" error from backend
    // TODO : Handle User phone numbers without a US country code

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
          className="m-3 place-self-center rounded-lg border-4 border-gray-200 bg-gray-200 p-0.5 text-black"
          name="email"
          placeholder="Email"
          required
        />
        <input
          type="password"
          className="m-3 place-self-center rounded-lg border-4 border-gray-200 bg-gray-200 p-0.5 text-black"
          name="password"
          placeholder="Password"
          required
        />
        <PhoneNumberInput />
        <button className="mx-auto w-40 bg-gray-200 text-black hover:bg-gray-300">
          Sign Up!
        </button>
        <p className="m-2 place-self-center text-xs">
          Already have an account?{" "}
          <a href="/login" className="text-blue-400">
            Login Here
          </a>
        </p>
      </div>
    </form>
  );
};
