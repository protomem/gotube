import React from "react";
import { Video } from "@/domain/entities";

import {
  AspectRatio,
  Avatar,
  Box,
  Card,
  CardFooter,
  Heading,
  Img,
  List,
  Text,
} from "@chakra-ui/react";
import { formatViews } from "@/lib/utils";

type VideoListItemProps = {
  video: Video;
};

const VideoListItem = ({ video }: VideoListItemProps) => {
  return (
    <Card flexDir="row">
      <AspectRatio width="340px" ratio={16 / 9}>
        <Img
          src="https://bit.ly/dan-abramov"
          alt="Dan Abramov"
          roundedLeft="md"
        />
      </AspectRatio>

      <CardFooter display="flex" flexDir="column">
        <Heading fontSize="lg">{video.title}</Heading>
        <Text>{formatViews(video.views) + " views"}</Text>
        <Box mt={2} display="flex" flexDir="row" alignItems="center" gap={2}>
          <Avatar
            name="Dan Abrahmov"
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
