import React from "react";
import AppBar from "../components/app-bar";
import SideBar from "../components/side-bar";
import {
  Box,
  Drawer,
  DrawerContent,
  DrawerOverlay,
  Flex,
  useDisclosure,
} from "@chakra-ui/react";

type Props = {
  children: React.ReactNode;
  hideSideBar?: boolean;
};

const MainLayout = ({ children, hideSideBar }: Props) => {
  hideSideBar = hideSideBar || false;

  const { isOpen, onOpen, onClose } = useDisclosure();
  const handleSwtchSideBar = () => {
    isOpen ? onClose() : onOpen();
  };

  return (
    <Flex direction="column" h="100dvh">
      <AppBar switchSideBar={handleSwtchSideBar} />

      <Flex direction="row" h="full">
        {!hideSideBar ? (
          <SideBar type={isOpen ? "minimal" : "expanded"} />
        ) : (
          <Drawer isOpen={isOpen} onClose={onClose} placement="left">
            <DrawerOverlay />
            <DrawerContent>
              <SideBar />
            </DrawerContent>
          </Drawer>
        )}

        <Box flex="1">{children}</Box>
      </Flex>
    </Flex>
  );
};

export default MainLayout;
