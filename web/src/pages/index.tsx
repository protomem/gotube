import { useSearchParams } from "next/navigation";

import AppBar from "@/components/app-bar";
import HomeNavMenu, { HomeNavMenuItemLabel } from "@/components/home-nav-menu";
import MainLayout from "@/components/layouts/main-layout";
import SideBar from "@/components/side-bar";
import { Box } from "@chakra-ui/react";

export default function Home() {
  const searchParams = useSearchParams();

  let selectedNavItem = HomeNavMenuItemLabel.New;
  if (searchParams.has("nav")) {
    switch (searchParams.get("nav")) {
      case HomeNavMenuItemLabel.New:
        selectedNavItem = HomeNavMenuItemLabel.New;
        break;
      case HomeNavMenuItemLabel.Popular:
        selectedNavItem = HomeNavMenuItemLabel.Popular;
        break;
      case HomeNavMenuItemLabel.Subscriptions:
        selectedNavItem = HomeNavMenuItemLabel.Subscriptions;
      default:
        break;
    }
  }

  return (
    <MainLayout
      appbar={<AppBar />}
      sidebar={
        <SideBar
          navmenu={
            <HomeNavMenu selectedItem={selectedNavItem} withSubscriptions />
          }
        />
      }
    >
      <Box>{"Body"}</Box>
    </MainLayout>
  );
}
