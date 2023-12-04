import { useRouter } from "next/router";
import { useQuery } from "@tanstack/react-query";
import { videoService } from "@/domain/video.service";
import dynamic from "next/dynamic";
import { repeat } from "@/lib/utils";
import { comments } from "@/domain/fixtures/comments";
import { Video } from "@/domain/entities";

import AppBar from "@/components/app-bar";
import MainLayout from "@/components/layouts/main-layout";
import { Box, Divider, Heading } from "@chakra-ui/react";
import CommentList from "@/components/comment-list";
import VideoHeader from "@/components/video-header";
import RatingButtons from "@/components/rating-buttons";
import SubscribeButton from "@/components/subscribe-button";

const DynamicVideoPlayer = dynamic(() => import("@/components/video-player"), {
  ssr: false,
});

export default function Watch() {
  const router = useRouter();
  const { videoId } = router.query;

  const { data: video } = useQuery({
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

        <Divider my={5} />

        <VideoHeader
          video={video as Video}
          buttonSubscribe={
            <SubscribeButton onSubscribe={() => {}} onUnsubscribe={() => {}} />
          }
          buttonRatings={<RatingButtons />}
        />

        <Divider my={5} />

        <Box mx={20} display="flex" flexDirection="column" gap={5}>
          <Heading fontSize="lg" fontWeight="bold">
            {"2132 Comments"}
          </Heading>

          <CommentList comments={repeat(comments, 30)} />
        </Box>
      </Box>
    </MainLayout>
  );
}
