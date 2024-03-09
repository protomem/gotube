import { FC, PropsWithChildren } from "react";
import AppBar from "./AppBar";
import SideBar from "./SideBar";
import { Box } from "@chakra-ui/react";

type Props = {
  searchTerm?: string;
};

const MainLayout: FC<PropsWithChildren<Props>> = ({ children, searchTerm }) => {
  return (
    <Box>
      <AppBar searchTerm={searchTerm} />

      <SideBar />

      <Box>{children}</Box>
    </Box>
  );
};

export default MainLayout;
