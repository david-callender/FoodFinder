import Image from "next/image";

import pioneer from "./pioneer.png";

import type { FC } from "react";

export const MainBackground: FC = () => {
  {
    return (
      <Image
        src={pioneer} // ⬅️ your background image path
        alt="Background"
        fill
        priority
        className="-z-10 object-cover object-center"
      />
    );
  }
};
