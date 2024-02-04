import {
  Button,
  Flex,
  Heading,
  IconButton,
  Input,
  InputGroup,
  InputRightAddon,
  Spacer,
} from "@chakra-ui/react";
import { FaSistrix } from "react-icons/fa6";

export default function AppBar() {
  return (
    <Flex
      h="100%"
      mx="4"
      direction="row"
      justifyItems="center"
      alignItems="center"
      gap="4"
    >
      <Heading size="lg">GoTube</Heading>

      <Spacer />

      <InputGroup maxW="2xl">
        <Input rounded="full" placeholder="Search ..." textAlign="center" />
        <InputRightAddon rounded="full" px="2">
          <IconButton
            aria-label="Search"
            rounded="full"
            variant="ghost"
            icon={<FaSistrix />}
          />
        </InputRightAddon>
      </InputGroup>

      <Spacer />

      <Button>Login</Button>
    </Flex>
  );
}
