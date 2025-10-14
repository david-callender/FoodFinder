import { MealSearch } from "@/components/MealSearch/MealSearch";
import { getMenu } from "@/db/getMenu";

import { Menu } from "../../components/Menu/Menu";

import type { FC, FormEvent } from "react";

// used for onSubmit behavior for MealQuery component. should match handleSubmit function
export type HandleMealQueryFunction = (
  event: FormEvent<HTMLFormElement>
) => void;

export const Menu_Page: FC = async () => {
  // need to store this here so Menu Component can acces the future response object

  const items = await getMenu(new Date(), "breakfast", "");
  return (
    <>
      <MealSearch />
      {/* TODO : update MenuItems from static -> response object from backend db
          MenuItems should be a state variable which is updated in handlSubmit */}
      <Menu items={items} />
    </>
  );
};

export default Menu_Page;
