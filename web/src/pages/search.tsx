import { useSearchParams } from "next/navigation";
import { videoService } from "@/domain/video.service";
import { useInfiniteQuery } from "@tanstack/react-query";

import AppBar from "@/components/app-bar";
import HomeNavMenu from "@/components/home-nav-menu";
import MainLayout from "@/components/layouts/main-layout";
import SideBar from "@/components/side-bar";
import VideoList from "@/components/video-list";
import { Box } from "@chakra-ui/react";

export default function Search() {
  const searchParams = useSearchParams();

  let searchQuery = "";
  if (searchParams.has("q")) {
    searchQuery = searchParams.get("q") || "";
  }

  const { data, fetchNextPage, hasNextPage } = useInfiniteQuery({
    queryKey: ["videos", { searchQuery }],
    queryFn: async ({ pageParam }) => {
      const limit = 6;
      return await videoService.searchVideos(searchQuery, {
        limit,
        offset: pageParam * limit,
      });
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage, _, lastPageParam) => {
      if (lastPage.data.videos.length !== 6) {
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
      appbar={<AppBar searchQuery={searchQuery} />}
      sidebar={<SideBar navmenu={<HomeNavMenu />} />}
    >
      <Box w="auto" height="full" px={5} py={3} sx={{ overflowY: "auto" }}>
        <VideoList
          videos={data ?? []}
          onLast={() => {
            if (hasNextPage) fetchNextPage();
          }}
        />
      </Box>
    </MainLayout>
  );
}
