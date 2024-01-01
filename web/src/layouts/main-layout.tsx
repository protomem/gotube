import React from "react";
import { Center } from "@chakra-ui/react";

type Props = {
  children: React.ReactNode;
};

const MainLayout = ({ children }: Props) => {
  return <Center pt="10">{children}</Center>;
};

export default MainLayout;
