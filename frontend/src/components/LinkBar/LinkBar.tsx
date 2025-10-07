import type { FC } from "react";

export const LinkBar: FC = () => {
  return (
    <div className="m-2 flex flex-row">
      <div className="grow text-center">
        <a href="/login">login</a>
      </div>
      <div className="grow text-center">
        <a href="/signup">signup</a>
      </div>
    </div>
  );
};
