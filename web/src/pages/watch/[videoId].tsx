import { useRouter } from "next/router";
import { useQuery } from "@tanstack/react-query";
import { videoService } from "@/domain/video.service";
import dynamic from "next/dynamic";

import AppBar from "@/components/app-bar";
import MainLayout from "@/components/layouts/main-layout";
import { Box, Center, Heading } from "@chakra-ui/react";

const DynamicVideoPlayer = dynamic(() => import("@/components/video-player"), {
  ssr: false,
});

export default function Watch() {
  const router = useRouter();
  const { videoId } = router.query;

  const { data } = useQuery({
    queryKey: ["video", videoId],
    queryFn: async () => {
      return await videoService.getVideo(videoId as string);
    },
    select: (data) => data.data.video,
    enabled: !!videoId,
  });

  return (
    <MainLayout appbar=<AppBar />>
      <Box w="auto" height="full" sx={{ overflowY: "auto" }}>
        <DynamicVideoPlayer
          src={
            "https://storage.googleapis.com/gtv-videos-bucket/sample/ElephantsDream.mp4"
          }
        />

        <Box>
          <Heading>{"Some Title"}</Heading>
        </Box>

        <Box>{"Comments ..."}</Box>
      </Box>
    </MainLayout>
  );
}
