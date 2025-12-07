import { FC } from "react";
import Image from "next/image";

import GithubLogo from "../components/res/github-mark-white.svg";

export const GithubLink: FC = () => {
   return (
            <a href="https://github.com/david-callender/FoodFinder"><Image
            src={GithubLogo}
            width={20}
            height={20}
            alt="Gopher Grub Logo"
            className="ml-50"
            />
        </a>
    )
}