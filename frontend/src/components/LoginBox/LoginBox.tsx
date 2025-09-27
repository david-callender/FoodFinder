"use client";

import { useRouter } from "next/navigation";

import type { FC, FormEvent } from "react";

export const LoginBox: FC = () => {
  // Description: email/password text field and login button w/ redirect
  const textInputClass =
    "bg-gray-200 place-self-center border-4 p-0.5 m-2 rounded-lg text-black";

  // use for redirects
  const router = useRouter();

  function handleSubmit(event: FormEvent<HTMLFormElement>): void {
    // prevents refresh of page
    event.preventDefault();

    // for use later
    // const formData = new FormData(event.currentTarget);
    // const email = formData.get('email');
    // const password = formData.get('password');

    // handle authentication here

    // redirect to home page
    router.push("/");
  }

  return (
    <form onSubmit={handleSubmit}>
      <div className="m-20 flex flex-auto flex-col p-5">
        <input
          type="email"
          name="email"
          placeholder="Email"
          className={textInputClass}
          required
        />
        <input
          type="password"
          name="password"
          placeholder="Password"
          className={textInputClass}
          required
        />
        <button className="mx-auto w-40 bg-gray-200 text-black hover:bg-gray-300">
          login
        </button>
      </div>
    </form>
  );
};
