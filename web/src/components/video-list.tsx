import React from "react";
import { Video } from "@/domain/entities";
import { formatDate, formatViews } from "@/lib/utils";

import {
  AspectRatio,
  Avatar,
  Box,
  Card,
  CardFooter,
  Divider,
  Heading,
  Img,
  List,
  Text,
} from "@chakra-ui/react";

type VideoListItemProps = {
  video: Video;
};

const VideoListItem = ({ video }: VideoListItemProps) => {
  return (
    <Card flexDir="row">
      <AspectRatio width="340px" ratio={16 / 9}>
        <Img src={video.thumbnailPath} alt={video.title} roundedLeft="md" />
      </AspectRatio>

      <Box>
        <Divider orientation="vertical" />
      </Box>

      <CardFooter display="flex" flexDir="column">
        <Heading fontSize="lg">{video.title}</Heading>
        <Text>{`${formatDate(video.createdAt)} • ${formatViews(
          video.views,
        )} views`}</Text>
        <Box mt={2} display="flex" flexDir="row" alignItems="center" gap={2}>
          <Avatar
            src={video.author.avatarPath}
            name={video.author.nickname}
            size="md"
            w="40px"
            h="40px"
            rounded="full"
          />
          <Text fontSize="lg">{video.author.nickname}</Text>
        </Box>
        <Text mt={4}>{video.description}</Text>
      </CardFooter>
    </Card>
  );
};

type VideoListProps = {
  videos: Video[];
};

const VideoList = ({ videos }: VideoListProps) => {
  return (
    <List spacing={3}>
      {videos.map((video) => (
        <VideoListItem key={video.id} video={video} />
      ))}
    </List>
  );
};

export default VideoList;
