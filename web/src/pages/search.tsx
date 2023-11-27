import { useState } from "react";
import { useSearchParams } from "next/navigation";
import { GetVideosParams, videoService } from "@/domain/video.service";
import { useQuery } from "@tanstack/react-query";

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

  const [videoParams, setVideoParams] = useState<GetVideosParams>({
    limit: 9,
    offset: 0,
  });

  const { data } = useQuery({
    queryKey: ["videos", { type: "search", query: searchQuery }],
    queryFn: () => videoService.searchVideos(searchQuery, { ...videoParams }),
    select: (data) => data.data.videos,
  });

  return (
    <MainLayout
      appbar={<AppBar searchQuery={searchQuery} />}
      sidebar={<SideBar navmenu={<HomeNavMenu />} />}
    >
      <Box w="auto" height="full" px={5} py={3} sx={{ overflowY: "auto" }}>
        <VideoList videos={data ?? []} />
      </Box>
    </MainLayout>
  );
}
