import { useNavigate } from "react-router-dom";
import AuthForm from "../components/auth-form";
import { FaBackward } from "react-icons/fa";
import { Box, Center, Flex, IconButton } from "@chakra-ui/react";

const AuthPage = () => {
  const nav = useNavigate();
  const handleBackClick = () => {
    nav("/", { replace: true });
  };

  return (
    <Flex h="100dvh" direction="column" gap="8rem">
      <Box w="full" paddingTop="10" paddingLeft="10">
        <IconButton
          aria-label="back"
          icon={<FaBackward />}
          variant="outline"
          size="lg"
          onClick={handleBackClick}
        />
      </Box>

      <Center>
        <AuthForm />
      </Center>
    </Flex>
  );
};

export default AuthPage;
