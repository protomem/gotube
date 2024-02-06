import { ReactNode } from "react";
import { useSideBarStore } from "@/domain/stores/side-bar";
import { Box, Drawer, DrawerBody, DrawerOverlay, Flex } from "@chakra-ui/react";

interface Props {
  appBar: ReactNode;
  sideBar?: ReactNode;
  hideSideBar?: boolean;
  children: ReactNode;
}

export default function MainLayout({
  appBar,
  sideBar,
  hideSideBar = false,
  children,
}: Props) {
  const [isOpen, close] = useSideBarStore((state) => [
    state.isOpen,
    state.close,
  ]);

  return (
    <Flex w="100%" h="100svh" direction="column">
      <Box h="7svh">{appBar}</Box>

      <Flex w="100%" h="93svh" direction="row">
        {sideBar && !hideSideBar ? (
          <Box w={isOpen ? "2xs" : "4rem"}>{sideBar}</Box>
        ) : (
          <Drawer isOpen={isOpen} placement="left" onClose={close}>
            <DrawerOverlay />
            <DrawerBody>cdscd</DrawerBody>
          </Drawer>
        )}

        <Box w="100%" overflowY="auto">
          {children}
        </Box>
      </Flex>
    </Flex>
  );
}
