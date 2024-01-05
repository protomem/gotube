import Logo from "./logo";
import Searcher from "./searcher";
import LoginButton from "./login-button";
import { Box, HStack, Spacer } from "@chakra-ui/react";

type Props = {
  switchSideBar: () => void;
};

const AppBar = ({ switchSideBar }: Props) => {
  return (
    <HStack h="16" px="6">
      <Logo switchSideBar={switchSideBar} />

      <Spacer />

      <Box>
        <Searcher />
      </Box>

      <Spacer />

      <Box>
        <LoginButton />
      </Box>
    </HStack>
  );
};

export default AppBar;
