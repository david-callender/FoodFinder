"use client";

import { useEffect, useState } from "react";

import type { FC } from "react";

export const MainPageWelcome: FC = () => {
  const [displayName, setDisplayName] = useState("Student");

  useEffect(() => {
    const displayName = localStorage.getItem("displayName");
    setDisplayName(displayName ?? "Student");
  }, []);

  return (
    <>
      {/* /menu Naviagtor */}
      <main className="relative z-10 flex flex-1 flex-col items-center justify-center px-4 text-center">
        <h1 className="mb-4 text-5xl font-bold text-white drop-shadow-lg">
          Gopher Grub
        </h1>
        <p className="text-white">
          Welcome, {displayName}!
        </p>
        <p className="mb-8 max-w-lg text-xl text-gray-100">
          Get notified about your favorite dining hall meals
        </p>
        <a
          href="/menu"
          className="rounded-xl bg-red-900 px-4 py-2 font-semibold text-white shadow drop-shadow-2xl transition hover:bg-red-700"
        >
          View Menu
        </a>
      </main>
    </>
  );
};
