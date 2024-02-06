import { useSearchParams } from "next/navigation";
import { useInfiniteQuery } from "@tanstack/react-query";
import videoService from "@/domain/video.service";
import { Box } from "@chakra-ui/react";
import MainLayout from "@/layouts/main-layout";
import AppBar from "@/widgets/app-bar";
import SideBar from "@/widgets/side-bar";
import VideoGrid from "@/components/video-grid";

const DEFAULT_LIMIT = 12;

export default function Home() {
  const searchParams = useSearchParams();

  const navMentItem = searchParams.get("videos") || "home";

  const fetchVideos = async (offset: number) => {
    switch (navMentItem) {
      case "home":
        return await videoService.findLatestVideos({
          offset,
          limit: DEFAULT_LIMIT,
        });
      case "trends":
        return await videoService.findTrendingVideos({
          offset,
          limit: DEFAULT_LIMIT,
        });
      default:
        throw new Error("Invalid navMentItem");
    }
  };

  const { data } = useInfiniteQuery({
    queryKey: ["videos", navMentItem],
    queryFn: async (params) => await fetchVideos(params.pageParam),
    initialPageParam: 0,
    getNextPageParam: (resp, _, offset) => {
      if (resp.data.videos.length < DEFAULT_LIMIT) return undefined;
      return offset + DEFAULT_LIMIT;
    },
    select: (data) => {
      return data.pages.flatMap((page) => page.data.videos);
    },
  });

  return (
    <MainLayout
      appBar={<AppBar />}
      sideBar={<SideBar navMenuItemSelected={navMentItem} />}
    >
      <Box px="4" mt="4">
        <VideoGrid videos={data} />
      </Box>
    </MainLayout>
  );
}
