import { LinkBar } from "@/components/LinkBar/LinkBar";
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
      <LinkBar />
      <MenuManager />
    </>
  );
};

export default Menu_Page;
