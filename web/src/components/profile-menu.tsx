import { useAuth } from "../providers/auth-provider";
import { User } from "../domain/entities";
import { resolveAddr } from "../domain/api.client";
import {
  Avatar,
  Button,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
} from "@chakra-ui/react";

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
        rightIcon={
          <Avatar
            name={user.nickname}
            src={resolveAddr(user.avatarPath)}
            size="sm"
          />
        }
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
