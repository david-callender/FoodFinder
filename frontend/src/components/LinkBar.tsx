"use client";

import Image from "next/image";
import Link from "next/link";
import { useEffect, useState } from "react";

import GopherGrubLogo from "./res/GopherGrubLogo.png";

import type { FC } from "react";

export const LinkBar: FC = () => {
  const [displayName, setDisplayName] = useState("Student");

  useEffect(() => {
    const displayName = localStorage.getItem("displayName");
    setDisplayName(displayName ?? "Student");
  }, []);

  return (
    <>
      <nav className="w-full border-b border-gray-200 bg-black/20">
        <div className="grid grid-cols-3 justify-center p-2">
          {/* Logo / Site Name */}
          <Link
            href="/"
            className="justify-self-start text-2xl font-bold text-white"
          >
            Gopher Grub
          </Link>

          <Image
            src={GopherGrubLogo}
            width={50}
            height={50}
            alt="Gopher Grub Logo"
            className="justify-self-center"
          />

          <div className="self-center justify-self-end">
            <span className="mr-2 text-white">{displayName}</span>
            <Link
              href="/login"
              className="mr-2 rounded-xl bg-red-900 px-4 py-2 font-semibold text-white shadow transition hover:bg-red-700"
            >
              Login
            </Link>
            <Link
              href="/signup"
              className="rounded-xl bg-red-900 px-4 py-2 font-semibold text-white shadow transition hover:bg-red-700"
            >
              Sign Up
            </Link>
          </div>
        </div>
      </nav>
    </>
  );
};
