import NavMenu from "@/components/nav-menu";
import { VStack } from "@chakra-ui/react";

interface Props {
  navMenuItemSelected?: string;
}

export default function SideBar({ navMenuItemSelected }: Props) {
  return (
    <VStack pl="2">
      <NavMenu labelSelected={navMenuItemSelected} />
    </VStack>
  );
}
