import { LinkBar } from "@/components/LinkBar/LinkBar";

import { SignUpBox } from "../../components/SignUpBox/SignUpBox";

import type { FC } from "react";

const Login: FC = () => {
  return (
    <>
      <LinkBar />
      <SignUpBox />
    </>
  );
};

export default Login;
