import { useState, ChangeEvent } from "react";
import { SearchIcon } from "@chakra-ui/icons";
import {
  IconButton,
  Input,
  InputGroup,
  InputRightAddon,
} from "@chakra-ui/react";

type Props = {
  defaultQuery?: string;
};

const Searcher = ({ defaultQuery }: Props) => {
  const [query, setQuery] = useState(defaultQuery);

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    setQuery(e.target.value);
  };

  return (
    <InputGroup>
      <Input
        placeholder="Search..."
        _placeholder={{ color: "gray.200" }}
        value={query}
        onChange={handleChange}
        w={{ sm: "xs", md: "md", lg: "lg", xl: "2xl" }}
        rounded="full"
        sx={{ textAlign: "center" }}
      />
      ;
      <InputRightAddon rounded="full" pl="2" pr="4">
        <IconButton
          aria-label="search"
          icon={<SearchIcon />}
          variant="ghost"
          rounded="full"
        />
      </InputRightAddon>
    </InputGroup>
  );
};

export default Searcher;
