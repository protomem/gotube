import NextLink from "next/link";
import { Button } from "@chakra-ui/react";

export default function ProfileMenu() {
  return (
    <Button as={NextLink} href="/auth/sign-in">
      login
    </Button>
  );
}
