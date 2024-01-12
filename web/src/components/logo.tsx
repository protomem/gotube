import { FaBars } from "react-icons/fa";
import { HStack, IconButton } from "@chakra-ui/react";
import Header from "./header";

type Props = {
  switchSideBar?: () => void;
};

const Logo = ({ switchSideBar }: Props) => {
  return (
    <HStack>
      <IconButton
        aria-label="side menu"
        icon={<FaBars />}
        variant="ghost"
        onClick={switchSideBar}
      />

      <Header />
    </HStack>
  );
};

export default Logo;
