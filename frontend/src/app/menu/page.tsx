import { LinkBar } from "@/components/LinkBar/LinkBar";
import { MainBackground } from "@/components/MainBackground/MainBackground";
import { MenuManager } from "@/components/MenuManager/MenuManager";

import type { FC, FormEvent } from "react";

// used for onSubmit behavior for MealQuery component. should match handleSubmit function
// TODO [misc.] : clean this type declaration up
export type HandleMealQueryFunction = (
  event: FormEvent<HTMLFormElement>
) => void;

export const Menu_Page: FC = () => {
  return (
    <>
      <div className="fixed inset-0 -z-10">
        {/* Optional dark overlay to make text pop */}
        <MainBackground />
        <div className="absolute inset-0 bg-black/40 backdrop-blur-sm"></div>
      </div>
      <LinkBar />
      <MenuManager />
    </>
  );
};

export default Menu_Page;
