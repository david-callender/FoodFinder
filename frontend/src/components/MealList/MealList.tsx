import type { FC } from "react";
import { MenuItem } from "../Menu/Menu";
import { LikeButton } from "../LikeButton/LikeButton";


type Props = {
    items: MenuItem[];
}

export const MealList: FC<Props> = ({ items }) => {
    {
        return (
            <>
            
                {items.map((item: MenuItem) => (
                    <tr key={item.id}>
                        <td className="p-5">{item.meal}</td>
                        <td className="p-5">
                            <LikeButton item={item} />
                        </td>
                    </tr>
                ))}

            </>
        )

    }
}

