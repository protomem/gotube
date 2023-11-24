import { useSearchParams } from "next/navigation";
import { videoService } from "@/domain/video.service";

import AppBar from "@/components/app-bar";
import HomeNavMenu, { HomeNavMenuItemLabel } from "@/components/home-nav-menu";
import MainLayout from "@/components/layouts/main-layout";
import SideBar from "@/components/side-bar";
import VideoGrid from "@/components/video-grid";
import { repeat } from "@/lib/utils";
import { version } from "react";
import { Box, Center } from "@chakra-ui/react";

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

  const videos = videoService.getVideos();

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
      <Box w="auto" height="full" px={5} py={3} sx={{ overflowY: "auto" }}>
        <VideoGrid videos={repeat(videos, 14)} />
      </Box>
    </MainLayout>
  );
}
