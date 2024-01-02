import {} from "react";
import Logo from "./logo";
import { DragHandleIcon } from "@chakra-ui/icons";
import { Box, HStack, IconButton, Spacer } from "@chakra-ui/react";

type Props = {
  switchSideBar: () => void;
};

const AppBar = ({ switchSideBar }: Props) => {
  return (
    <HStack h="16" px="4">
      <HStack>
        <IconButton
          aria-label="side menu"
          icon={<DragHandleIcon />}
          variant="ghost"
          onClick={switchSideBar}
        />

        <Logo />
      </HStack>

      <Spacer />

      <Box>Search</Box>

      <Spacer />

      <Box>Login Button</Box>
    </HStack>
  );
};

export default AppBar;
