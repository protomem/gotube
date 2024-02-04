import NextLink from "next/link";
import { FaBars } from "react-icons/fa6";
import { Box, Heading, IconButton } from "@chakra-ui/react";

export default function Title() {
  return (
    <Box display="flex" flexDir="row" gap="2" justifyContent="center" alignItems="center">
      <IconButton
        aria-label="side-bar toggle"
        icon={<FaBars />}
        variant="ghost"
        onClick={() => { }}
      />

      <Heading as={NextLink} href="/" size="lg">
        GoTube
      </Heading>
    </Box>
  );
}
