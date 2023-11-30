import { useSearchParams } from "next/navigation";
import { videoService } from "@/domain/video.service";
import { useInfiniteQuery } from "@tanstack/react-query";

import AppBar from "@/components/app-bar";
import HomeNavMenu, { HomeNavMenuItemLabel } from "@/components/home-nav-menu";
import MainLayout from "@/components/layouts/main-layout";
import SideBar from "@/components/side-bar";
import VideoGrid from "@/components/video-grid";
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

  const { data, fetchNextPage, hasNextPage } = useInfiniteQuery({
    queryKey: ["videos"],
    queryFn: async ({ pageParam }) => {
      const limit = 9;
      switch (selectedNavItem) {
        case HomeNavMenuItemLabel.New:
          return await videoService.getNewVideos({
            limit,
            offset: pageParam * limit,
          });
        case HomeNavMenuItemLabel.Popular:
          return await videoService.getPopularVideos({
            limit,
            offset: pageParam * limit,
          });
        default:
          return await videoService.getNewVideos({
            limit,
            offset: pageParam * limit,
          });
      }
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage, _, lastPageParam) => {
      if (lastPage.data.videos.length !== 9) {
        return undefined;
      }
      return lastPageParam + 1;
    },
    select: (data) => {
      return data.pages.flatMap((page) => page.data.videos);
    },
  });

  return (
    <MainLayout
      appbar={<AppBar />}
      sidebar={
        <SideBar navmenu={<HomeNavMenu selectedItem={selectedNavItem} />} />
      }
    >
      <Box w="auto" height="full" px={5} py={3} sx={{ overflowY: "auto" }}>
        <VideoGrid
          videos={data !== undefined ? data : []}
          onLast={() => {
            if (hasNextPage) fetchNextPage();
          }}
        />
      </Box>
    </MainLayout>
  );
}
