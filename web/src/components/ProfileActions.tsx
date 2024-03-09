import { FC } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@chakra-ui/react";

const ProfileActions: FC = () => {
  const nav = useNavigate();

  const handleClick = () => {
    nav("/auth/sign-in");
  };

  return <Button onClick={handleClick}>Login</Button>;
};

export default ProfileActions;
