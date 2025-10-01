import type { FC } from "react";

export const LinkBar: FC = () => {
  return (
    <div className="m-2 mx-100 flex flex-row">
      <div className="basis-1/2 text-center">
        <a href="/login">login</a>
      </div>
      <div className="basis-1/2 text-center">
        <a href="/signup">signup</a>
      </div>
    </div>
  );
};
