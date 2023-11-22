import AppBar from "@/components/app-bar";
import MainLayout from "@/components/layouts/main-layout";
import SideBar from "@/components/side-bar";
import { Box } from "@chakra-ui/react";

export default function Home() {
  return (
    <MainLayout appbar=<AppBar /> sidebar=<SideBar />>
      <Box>{"Body"}</Box>
    </MainLayout>
  );
}
