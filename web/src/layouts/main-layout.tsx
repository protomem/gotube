import { ReactNode } from "react";
import { Box, Flex } from "@chakra-ui/react";

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
  return (
    <Flex w="100%" h="100svh" direction="column">
      <Box h="7svh">{appBar}</Box>

      <Flex w="100%" h="93svh" direction="row">
        {sideBar && !hideSideBar && <Box w="2xs">{sideBar}</Box>}

        <Box w="100%" overflowY="scroll">
          {children}
        </Box>
      </Flex>
    </Flex>
  );
}
