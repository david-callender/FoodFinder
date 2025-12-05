"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";

import { logout } from "@/db/logout";

import type { FC } from "react";

export const Logout: FC = () => {
  const router = useRouter();

  useEffect(() => {
    void (async (): Promise<void> => {
      const result = await logout();

      if (result.ok) {
        router.push("/");
      } else {
        alert("Error logging out: " + result.err);
      }
    })();
  }, [router]);

  return <></>;
};
