import { FC } from "react";
import useSearch from "../hooks/search";
import { FaMagnifyingGlass } from "react-icons/fa6";
import {
  Box,
  FormControl,
  IconButton,
  Input,
  InputGroup,
  InputRightAddon,
} from "@chakra-ui/react";

type Props = {
  defaultValue?: string;
};

const Searcher: FC<Props> = ({ defaultValue }) => {
  const { handleSubmit, inputProps } = useSearch(defaultValue);

  return (
    <Box maxW="3xl" minW="md" w="full">
      <form onSubmit={handleSubmit}>
        <FormControl>
          <InputGroup>
            <Input
              placeholder="Search ..."
              borderLeftRadius="full"
              textAlign="center"
              {...inputProps}
            />
            <InputRightAddon p="2" borderRightRadius="full">
              <IconButton
                aria-label="Search"
                icon={<FaMagnifyingGlass />}
                variant="ghost"
                borderRadius="full"
                type="submit"
              />
            </InputRightAddon>
          </InputGroup>
        </FormControl>
      </form>
    </Box>
  );
};

export default Searcher;
