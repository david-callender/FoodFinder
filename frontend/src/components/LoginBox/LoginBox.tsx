"use client";

import { useRouter } from "next/navigation";

import type { FC, FormEvent } from "react";
import { json } from "stream/consumers";

export const LoginBox: FC = () => {
  // Description: email/password text field and login button w/ redirect
  const textInputClass =
    "bg-gray-200 place-self-center border-4 p-0.5 m-2 rounded-lg text-black";

  // use for redirects
  const router = useRouter();

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    // prevents refresh of page
    event.preventDefault();

    // pulling form data
    const formData = new FormData(event.currentTarget);
    const email = formData.get('email');
    const password = formData.get('password');

    // TODO: input validation

    // handling authentication

    // 9/27: endpoint currently assumes all users are valid.
    // At some point, when validation is implemented,
    // update this to a redirect (?) for registration
    // if login fails

    // how to get refresh token?

    // FIX: Static URL
    // request to login endpoint
    const response = await fetch('http://localhost:8080/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: email, password: password }),
    })

    // redirect to home page
    if (response.ok){
     
        let userInfo = await response.json();
        // storing access token
        localStorage.setItem("access_token", userInfo.access_token);
        router.push("/");
    } else {
        console.log(await response.json())
    }
   
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
