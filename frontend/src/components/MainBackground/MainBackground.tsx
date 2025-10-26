import Image from "next/image";

import pioneer from "./pioneer.png";

import type { FC } from "react";

export const MainBackground: FC = () => {
  {
    return (
      <div className="fixed inset-0 -z-10">
        <Image
          src={pioneer}
          alt="Background"
          fill
          priority
          style={{ objectFit: "cover" }}
        />
      </div>
    )
    
  }
};
