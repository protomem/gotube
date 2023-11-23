import React from "react";

import { Flex, Spacer } from "@chakra-ui/react";
import Title from "@/components/title";
import ProfileMenu from "@/components/profile-menu";
import Searcher from "@/components/searcher";

const AppBar = () => {
  return (
    <Flex mx={4} my={2} alignItems="center">
      <Title />

      <Spacer />

      <Searcher />

      <Spacer />

      <ProfileMenu />
    </Flex>
  );
};

export default AppBar;
