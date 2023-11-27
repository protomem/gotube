import { useState } from "react";
import { useSearchParams } from "next/navigation";
import { GetVideosParams, videoService } from "@/domain/video.service";

import AppBar from "@/components/app-bar";
import HomeNavMenu, { HomeNavMenuItemLabel } from "@/components/home-nav-menu";
import MainLayout from "@/components/layouts/main-layout";
import SideBar from "@/components/side-bar";
import VideoGrid from "@/components/video-grid";
import { Box } from "@chakra-ui/react";
import { useQuery } from "@tanstack/react-query";

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

  const [videoParams, setVideoParams] = useState<GetVideosParams>({
    limit: 9,
    offset: 0,
  });

  const { data } = useQuery({
    queryKey: ["videos", { type: selectedNavItem }],
    queryFn: async () => {
      switch (selectedNavItem) {
        case HomeNavMenuItemLabel.New:
          return await videoService.getNewVideos({ ...videoParams });
        case HomeNavMenuItemLabel.Popular:
          return await videoService.getPopularVideos({ ...videoParams });
        default:
          return await videoService.getNewVideos({ ...videoParams });
      }
    },
    select: (data) => data.data.videos,
  });

  return (
    <MainLayout
      appbar={<AppBar />}
      sidebar={
        <SideBar navmenu={<HomeNavMenu selectedItem={selectedNavItem} />} />
      }
    >
      <Box w="auto" height="full" px={5} py={3} sx={{ overflowY: "auto" }}>
        <VideoGrid videos={data ?? []} />
      </Box>
    </MainLayout>
  );
}
