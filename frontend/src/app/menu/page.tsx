import { FinalMenu } from "@/components/FinalMenu/FinalMenu";

import type {FC, FormEvent} from "react";

// used for onSubmit behavior for MealQuery component. should match handleSubmit function
export type HandleMealQueryFunction = (
  event: FormEvent<HTMLFormElement>
) => void;

export const Menu_Page: FC = () => {
  
  return <FinalMenu />
  
};

export default Menu_Page;
