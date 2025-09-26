import type { FC } from "react";

// would it make more sense to keep components in respective page directory?

export const LoginBox: FC = () => {
  const textInputClass =
    "bg-gray-200 w-40 place-self-center border-4 border-transparent text-black m-2";
  return (
    <div className="grid border-4 text-black">
      <input type="text" className={textInputClass} />
      <input type="text" className={textInputClass} />
      <button>login</button>
    </div>
  );
};
