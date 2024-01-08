import { useAuth } from "../providers/auth-provider";
import { User } from "../domain/entities";
import { FaUser } from "react-icons/fa";
import { Button, Menu, MenuButton, MenuItem, MenuList } from "@chakra-ui/react";

type Props = {
  user: User;
};

const ProfileMenu = ({ user }: Props) => {
  const { logout } = useAuth();

  return (
    <Menu>
      <MenuButton
        as={Button}
        variant="ghost"
        size="lg"
        ml={2}
        rightIcon={<FaUser />}
        gap="2"
      >
        {user.nickname}
      </MenuButton>
      <MenuList p="0">
        <MenuItem>Profile</MenuItem>
        <MenuItem>Settings</MenuItem>
        <MenuItem>Studio</MenuItem>
        <MenuItem bg="red.500" onClick={logout}>
          Logout
        </MenuItem>
      </MenuList>
    </Menu>
  );
};

export default ProfileMenu;
