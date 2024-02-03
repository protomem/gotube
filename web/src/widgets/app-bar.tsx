import { Button, Flex, Heading, Input, Spacer } from "@chakra-ui/react";

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

      <Input
        maxW="2xl"
        rounded="full"
        placeholder="Search ..."
        textAlign="center"
      />

      <Spacer />

      <Button>Login</Button>
    </Flex>
  );
}
