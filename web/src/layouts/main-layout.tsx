import { ReactNode, useRef } from "react";
import { useSideBarState } from "../providers/side-bar-state-provider";
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
  useDimensions,
} from "@chakra-ui/react";

type Props = {
  children: ReactNode;
  hideSideBar?: boolean;
  selectedNavMenuItem?: NavMenuItem;
};

const MainLayout = ({ children, hideSideBar, selectedNavMenuItem }: Props) => {
  hideSideBar = hideSideBar || false;

  const appBarElement = useRef(null);
  const appBarDimensions = useDimensions(appBarElement);

  const { isOpen, onClose, onToggle: handleSwtchSideBar } = useSideBarState();

  return (
    <Flex direction="column" h="100dvh">
      <AppBar switchSideBar={handleSwtchSideBar} ref={appBarElement} />

      <Flex
        direction="row"
        h={`calc(100dvh - ${appBarDimensions?.borderBox.height}px)`}
      >
        {!hideSideBar ? (
          <Box overflow="auto" pr="2">
            <SideBar
              type={isOpen ? "minimal" : "expanded"}
              selectedNavMenuItem={selectedNavMenuItem}
            />
          </Box>
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
