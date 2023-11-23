import { useSearchParams } from "next/navigation";

import AppBar from "@/components/app-bar";
import HomeNavMenu from "@/components/home-nav-menu";
import MainLayout from "@/components/layouts/main-layout";
import SideBar from "@/components/side-bar";
import { Box } from "@chakra-ui/react";

export default function Search() {
  const searchParams = useSearchParams();

  let searchQuery = "";
  if (searchParams.has("q")) {
    searchQuery = searchParams.get("q") || "";
  }

  return (
    <MainLayout
      appbar={<AppBar searchQuery={searchQuery} />}
      sidebar={<SideBar navmenu={<HomeNavMenu withSubscriptions />} />}
    >
      <Box>{searchQuery}</Box>
    </MainLayout>
  );
}
