import React from "react";

import { VStack } from "@chakra-ui/react";

type SideBarProps = {
  navmenu?: React.ReactNode;
};

const SideBar = ({ navmenu }: SideBarProps) => {
  return (
    <VStack mx={2} align="stretch">
      {navmenu}
    </VStack>
  );
};

export default SideBar;
