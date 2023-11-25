import { useSearchParams } from "next/navigation";
import { videoService } from "@/domain/video.service";
import { repeat } from "@/lib/utils";

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

  const videos = videoService.getVideos();

  return (
    <MainLayout
      appbar={<AppBar searchQuery={searchQuery} />}
      sidebar={<SideBar navmenu={<HomeNavMenu withSubscriptions />} />}
    >
      <Box w="auto" height="full" px={5} py={3} sx={{ overflowY: "auto" }}>
        <VideoList videos={repeat(videos, 18)} />
      </Box>
    </MainLayout>
  );
}
