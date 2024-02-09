import { Flex, Spacer } from "@chakra-ui/react";
import Title from "@/components/title";
import Searcher from "@/components/searcher";
import ProfileMenu from "@/components/profile-menu";

interface Props {
  searchTerm?: string;
}

export default function AppBar({ searchTerm }: Props) {
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

      <Searcher defaultTerm={searchTerm} />

      <Spacer />

      <ProfileMenu />
    </Flex>
  );
}
