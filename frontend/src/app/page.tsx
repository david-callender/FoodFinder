"use client";

import { LinkBar } from "@/components/LinkBar/LinkBar";
import { MainBackground } from "@/components/MainBackground/MainBackground";

import type { FC } from "react";

const Home: FC = () => {
  return (
    <div className="relative flex min-h-screen flex-col overflow-hidden text-gray-900">
      <LinkBar />

      {/* Header */}
      <div className="absolute inset-0 -z-10">
        <MainBackground />
        {/* Optional dark overlay to make text pop */}
        <div className="absolute inset-0 bg-black/40 backdrop-blur-sm"></div>
      </div>

      {/* /menu Naviagtor */}
      <main className="relative z-10 flex flex-1 flex-col items-center justify-center px-4 text-center">
        <h1 className="mb-4 text-5xl font-bold text-white drop-shadow-lg">
          Gopher Grub
        </h1>
        <p className="mb-8 max-w-lg text-xl text-gray-100">
          Explore campus dining menus, find your favorite meals, and never miss
          taco night again.
        </p>
        <a
          href="/menu"
          className="rounded-xl bg-red-900 px-4 py-2 font-semibold text-white shadow drop-shadow-2xl transition hover:bg-red-700"
        >
          View Menu
        </a>
      </main>

      <footer className="relative z-10 border-t border-white/20 bg-black/30 py-4 text-center text-sm text-gray-200">
        © {new Date().getFullYear()} Gopher Grub · University of Minnesota
      </footer>
    </div>
  );
};

export default Home;
