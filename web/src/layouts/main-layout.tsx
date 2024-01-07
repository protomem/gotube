import { ReactNode } from "react";
import { useSideBarState } from "../app/side-bar-state-provider";
import AppBar from "../components/app-bar";
import SideBar from "../components/side-bar";
import Logo from "../components/logo";
import { NavMenuItem } from "../components/nav-menu";
import {
  Box,
  Drawer,
  DrawerContent,
  DrawerOverlay,
  Flex,
} from "@chakra-ui/react";

type Props = {
  children: ReactNode;
  hideSideBar?: boolean;
  selectedNavMenuItem?: NavMenuItem;
};

const MainLayout = ({ children, hideSideBar, selectedNavMenuItem }: Props) => {
  hideSideBar = hideSideBar || false;

  const { isOpen, onClose, onToggle: handleSwtchSideBar } = useSideBarState();

  return (
    <Flex direction="column" h="100dvh">
      <AppBar switchSideBar={handleSwtchSideBar} />

      <Flex direction="row" h="full">
        {!hideSideBar ? (
          <SideBar
            type={isOpen ? "minimal" : "expanded"}
            selectedNavMenuItem={selectedNavMenuItem}
          />
        ) : (
          <Drawer isOpen={isOpen} onClose={onClose} placement="left">
            <DrawerOverlay />
            <DrawerContent
              paddingTop="2"
              backgroundColor="gray.800"
              maxW="14rem"
            >
              <Box paddingLeft="5" paddingBottom="4">
                <Logo switchSideBar={handleSwtchSideBar} />
              </Box>

              <SideBar
                type="expanded"
                selectedNavMenuItem={selectedNavMenuItem}
              />
            </DrawerContent>
          </Drawer>
        )}

        <Box flex="1">{children}</Box>
      </Flex>
    </Flex>
  );
};

export default MainLayout;
