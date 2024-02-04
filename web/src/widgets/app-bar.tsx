import { Flex, Spacer } from "@chakra-ui/react";
import Title from "@/components/title";
import Searcher from "@/components/searcher";
import ProfileMenu from "@/components/profile-menu";

export default function AppBar() {
  return (
    <Flex
      h="100%"
      mx="4"
      direction="row"
      justifyItems="center"
      alignItems="center"
      gap="6"
    >
      <Title />

      <Spacer />

      <Searcher />

      <Spacer />

      <ProfileMenu />
    </Flex>
  );
}
