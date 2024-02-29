import { useAuthStore } from "@/domain/stores/auth";
import NextLink from "next/link";
import { FaRegUser } from "react-icons/fa6";
import {
  Button,
  IconButton,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
} from "@chakra-ui/react";

export default function ProfileMenu() {
  const [isAuthenticated, user, logout] = useAuthStore((state) => [
    state.isAuthenticated,
    state.user,
    state.logout,
  ]);

  return (
    <>
      {!(isAuthenticated() && user) ? (
        <Button as={NextLink} href="/auth/sign-in">
          login
        </Button>
      ) : (
        <Menu>
          <MenuButton
            as={IconButton}
            aria-label="Profile Menu"
            icon={<FaRegUser />}
            size="lg"
            rounded="full"
            variant="ghost"
          />

          <MenuList>
            <MenuItem as={NextLink} href={`/profile/${user.nickname}`}>
              Profile
            </MenuItem>
            <MenuItem as={NextLink} href={`/studio?section=profile`}>
              Settings
            </MenuItem>
            <MenuItem as={NextLink} href="/studio">
              Studio
            </MenuItem>
            <MenuItem
              as={Button}
              onClick={logout}
              colorScheme="red"
              justifyContent="start"
              variant="solid"
              rounded="none"
              size="sm"
            >
              Logout
            </MenuItem>
          </MenuList>
        </Menu>
      )}
    </>
  );
}
