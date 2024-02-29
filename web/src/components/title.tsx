import NextLink from "next/link";
import { FaBars } from "react-icons/fa6";
import { Box, Heading, IconButton } from "@chakra-ui/react";
import { useSideBarStore } from "@/domain/stores/side-bar";

export default function Title() {
  const [isOpen, open, close] = useSideBarStore((state) => [
    state.isOpen,
    state.open,
    state.close,
  ]);
  const toogle = () => {
    !isOpen ? open() : close();
  };

  return (
    <Box
      display="flex"
      flexDir="row"
      gap="2"
      justifyContent="center"
      alignItems="center"
    >
      <IconButton
        aria-label="side-bar toggle"
        icon={<FaBars />}
        variant="ghost"
        onClick={toogle}
      />

      <Heading as={NextLink} href="/" size="lg">
        GoTube
      </Heading>
    </Box>
  );
}
