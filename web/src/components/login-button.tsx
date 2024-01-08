import { useNavigate } from "react-router-dom";
import { Button } from "@chakra-ui/react";

const LoginButton = () => {
  const nav = useNavigate();
  const handleClick = () => {
    nav("/auth", { replace: true });
  };

  return (
    <Button onClick={handleClick} size="lg">
      login
    </Button>
  );
};

export default LoginButton;
