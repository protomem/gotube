import React, { useState } from "react";

import {
  Box,
  IconButton,
  Input,
  InputGroup,
  InputRightAddon,
} from "@chakra-ui/react";
import { SearchIcon } from "@chakra-ui/icons";
import { useRouter } from "next/navigation";
import { ROUTES, withQuery } from "@/lib/routes";

type SearcherProps = {
  defaultQuery?: string;
};

const Searcher = ({ defaultQuery }: SearcherProps) => {
  const router = useRouter();

  const [inputValue, setInputValue] = useState(defaultQuery ?? "");

  return (
    <Box w="350px">
      <InputGroup>
        <Input
          value={inputValue === "" ? undefined : inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter")
              router.push(withQuery(ROUTES.SEARCH, { q: inputValue }));
          }}
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
            onClick={() => {
              router.push(withQuery(ROUTES.SEARCH, { q: inputValue }));
            }}
          />
        </InputRightAddon>
      </InputGroup>
    </Box>
  );
};

export default Searcher;
