import { FC } from "react";
import { HStack } from "@chakra-ui/react";
import Title from "./Title";
import Searcher from "./Searcher";
import ProfileActions from "./ProfileActions";

type Props = {
  searchTerm?: string;
};

const AppBar: FC<Props> = ({ searchTerm }) => {
  return (
    <HStack
      justify="space-between"
      align="center"
      p="2"
      gap="6"
      bg="gray.700"
      sx={{ position: "sticky", top: 0 }}
    >
      <Title />

      <Searcher defaultValue={searchTerm} />

      <ProfileActions />
    </HStack>
  );
};

export default AppBar;
