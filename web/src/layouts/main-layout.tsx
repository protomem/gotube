import { ReactNode } from "react";
import { Box, Flex } from "@chakra-ui/react";

interface Props {
  appbar: ReactNode;
  children: ReactNode;
}

export default function MainLayout({ appbar, children }: Props) {
  return (
    <Flex w="100%" h="100svh" direction="column">
      <Box h="7svh">{appbar}</Box>

      <Box h="93svh" overflowY="auto">
        {children}
      </Box>
    </Flex>
  );
}
