import { useAuth } from "../providers/auth-provider";
import Logo from "./logo";
import Searcher from "./searcher";
import LoginButton from "./login-button";
import { Box, HStack, Spacer, forwardRef } from "@chakra-ui/react";
import ProfileMenu from "./profile-menu";

type Props = {
  switchSideBar: () => void;
};

const AppBar = forwardRef((props: Props, ref) => {
  const { switchSideBar } = props;
  const { isAuthenticated, currentUser } = useAuth();

  return (
    <HStack h="16" px="6" ref={ref}>
      <Logo switchSideBar={switchSideBar} />

      <Spacer />

      <Box>
        <Searcher />
      </Box>

      <Spacer />

      <Box>
        {isAuthenticated && currentUser ? (
          <ProfileMenu user={currentUser} />
        ) : (
          <LoginButton />
        )}
      </Box>
    </HStack>
  );
});

export default AppBar;
