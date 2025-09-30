import type { FC } from "react";

export const SignUpBox: FC = () => {
  return (
    <form>
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
          placeholder="Phone Number"
        />
        <button className="mx-auto w-40 bg-gray-200 text-black hover:bg-gray-300">
          Sign Up!
        </button>
      </div>
    </form>
  );
};
