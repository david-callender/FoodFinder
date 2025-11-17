import Link from "next/link";

import type { FC } from "react";

export const Header: FC = () => {
  return (
    <>
      {/* Header */}
      <header className="bg-maroon-700 py-8 text-center text-white shadow-md">
        <h1 className="text-5xl font-bold tracking-tight text-amber-800">
          <Link href="/">GopherGrub</Link>
        </h1>
        <p className="mt-2 text-lg text-amber-800 opacity-90">
          University of Minnesota Dining Made Simple
        </p>
      </header>
    </>
  );
};
