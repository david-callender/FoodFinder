"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";

import { logout } from "@/db/logout";

import type { FC } from "react";

export const Logout: FC = () => {
  const router = useRouter();

  useEffect(() => {
    void logout();
    router.push("/");
  }, [router]);

  return <></>;
};
