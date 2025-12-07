import { GithubLink } from "./GithubLink";

import type { FC } from "react";

export const Footer: FC = () => {
  return (
    <>
      <div className="grid grid-cols-3 content-center justify-center border-t border-gray-200 bg-black/30 p-2 py-4 text-center text-sm text-gray-200">
        <>
          <p></p>
        </>

        <footer className="relative z-10">
          <p>
            © {new Date().getFullYear()} Gopher Grub · University of Minnesota
          </p>
        </footer>

        <div className="mr-5 self-center justify-self-end">
          <GithubLink />
        </div>
      </div>
    </>
  );
};
