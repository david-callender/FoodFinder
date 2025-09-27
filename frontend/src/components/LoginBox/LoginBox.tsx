"use client";

import type { FC } from "react";

// would it make more sense to keep components in respective page directory?

export const LoginBox: FC = () => {
  const textInputClass =
    "bg-gray-200 place-self-center border-4 p-0.5 m-2 rounded-lg text-black";
  return (
    <form>
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
