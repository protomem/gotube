import { FC } from "react";
import { useSearchParams } from "react-router-dom";
import { Box, Heading, Text } from "@chakra-ui/react";
import MainLayout from "../components/MainLayout";

const Search: FC = () => {
  const [searchParams] = useSearchParams();
  const searchTerm = searchParams.get("term") || undefined;

  return (
    <MainLayout searchTerm={searchTerm}>
      <Box display="flex" flexDir="row" alignItems="end" gap={2}>
        <Heading>Search:</Heading>
        <Text as="i" fontSize="3xl" textDecoration="underline">
          {searchTerm}
        </Text>
      </Box>
    </MainLayout>
  );
};

export default Search;
