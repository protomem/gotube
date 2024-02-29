import { ReactNode } from "react";
import { Flex } from "@chakra-ui/react";

interface Props {
  children: ReactNode;
}

export default function SingleObjectLayout({ children }: Props) {
  return (
    <Flex align="center" justify="center" h="100svh">
      {children}
    </Flex>
  );
}
