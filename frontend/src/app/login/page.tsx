import { LinkBar } from "@/components/LinkBar/LinkBar";

import { LoginBox } from "../../components/LoginBox/LoginBox";

import type { FC } from "react";

const Login: FC = () => {
  return (
    <>
      <LinkBar />
      <LoginBox />
    </>
  );
};

export default Login;
