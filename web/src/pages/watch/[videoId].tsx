import { useRouter } from "next/router";
import { useQuery } from "@tanstack/react-query";
import { videoService } from "@/domain/video.service";
import { commentService } from "@/domain/comment.service";
import dynamic from "next/dynamic";

import AppBar from "@/components/app-bar";
import MainLayout from "@/components/layouts/main-layout";
import { Box, Divider, Heading } from "@chakra-ui/react";
import CommentList from "@/components/comment-list";
import VideoHeader from "@/components/video-header";
import RatingButtons from "@/components/rating-buttons";
import SubscribeButton from "@/components/subscribe-button";
import Error from "next/error";

const DynamicVideoPlayer = dynamic(() => import("@/components/video-player"), {
  ssr: false,
});

export default function Watch() {
  const router = useRouter();
  const { videoId } = router.query;

  const { data: video, isError } = useQuery({
    queryKey: ["video", videoId],
    queryFn: async () => {
      return await videoService.getVideo(videoId as string);
    },
    select: (data) => data.data.video,
    enabled: !!videoId,
  });

  const { data: comments } = useQuery({
    queryKey: ["comments", { video: videoId }],
    queryFn: async () => {
      return await commentService.getCommentsByVideo(videoId as string);
    },
    select: (data) => data.data.comments,
    enabled: !!videoId,
  });

  if (isError) {
    return <Error statusCode={404} />;
  }

  return (
    <MainLayout appbar=<AppBar />>
      <Box w="auto" height="full" sx={{ overflowY: "auto" }}>
        <DynamicVideoPlayer src={video?.videoPath || ""} />

        <Divider my={5} />

        {video && (
          <VideoHeader
            video={video}
            subscribers={1_000_000}
            buttonSubscribe={
              <SubscribeButton
                onSubscribe={() => {}}
                onUnsubscribe={() => {}}
              />
            }
            buttonRatings={<RatingButtons />}
          />
        )}

        <Divider my={5} />

        <Box mx={20} mb={10} display="flex" flexDirection="column" gap={5}>
          <Heading fontSize="lg" fontWeight="bold">
            {`${comments !== undefined ? comments.length : 0} Comments`}
          </Heading>

          <CommentList comments={comments || []} />
        </Box>
      </Box>
    </MainLayout>
  );
}
