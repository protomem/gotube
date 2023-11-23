import React from "react";

import {
  Box,
  IconButton,
  Input,
  InputGroup,
  InputRightAddon,
} from "@chakra-ui/react";
import { SearchIcon } from "@chakra-ui/icons";

const Searcher = () => {
  return (
    <Box w="350px">
      <InputGroup>
        <Input
          variant="filled"
          placeholder="Search..."
          rounded="full"
          sx={{ textAlign: "center" }}
        />
        <InputRightAddon rounded="full" p={1} pl={2}>
          <IconButton
            aria-label="search"
            m={0}
            rounded="full"
            icon={<SearchIcon />}
            variant="ghost"
          />
        </InputRightAddon>
      </InputGroup>
    </Box>
  );
};

export default Searcher;
