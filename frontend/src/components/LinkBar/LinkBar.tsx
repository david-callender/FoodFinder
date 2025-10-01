import type { FC } from "react";

export const LinkBar: FC = () => {
  return (
    <div className="m-2 flex flex-row">
      <div className="text-center grow">
        <a href="/login">login</a>
      </div>
      <div className="text-center grow">
        <a href="/signup">signup</a>
      </div>
    </div>
  );
};
