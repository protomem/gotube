import { useSearchParams } from "next/navigation";
import { useInfiniteQuery } from "@tanstack/react-query";
import videoService from "@/domain/video.service";
import MainLayout from "@/layouts/main-layout";
import AppBar from "@/widgets/app-bar";
import SideBar from "@/widgets/side-bar";
import { Box, Divider, HStack, Text } from "@chakra-ui/react";
import VideoList from "@/components/video-list";

const DEFAULT_LIMIT = 9;

export default function Search() {
  const searchParams = useSearchParams();

  const searchTerm = searchParams.get("term") || "";

  const fetchVideos = async (offset: number) => {
    return await videoService.searchVideos(searchTerm, {
      offset,
      limit: DEFAULT_LIMIT,
    });
  };

  const { data } = useInfiniteQuery({
    queryKey: ["videos", "search", searchTerm],
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
      appBar={<AppBar searchTerm={searchTerm} />}
      sideBar={<SideBar />}
    >
      <Box px="4">
        <HStack align="end" pl="2">
          <Text fontSize="lg" fontWeight="semibold">{`Search results:`}</Text>
          <Text fontSize="md" fontStyle="italic">{`${searchTerm}`}</Text>
        </HStack>

        <Divider mt="2" mb="4" />

        <VideoList videos={data} />
      </Box>
    </MainLayout>
  );
}
