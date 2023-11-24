import React from "react";
import { Video } from "@/domain/entities";

import {
  Card,
  CardBody,
  CardFooter,
  CardHeader,
  SimpleGrid,
} from "@chakra-ui/react";

type VideoGridItemProps = {
  video: Video;
};

const VideoGridItem = ({ video }: VideoGridItemProps) => {
  return (
    <Card minW="300px" w="auto" maxW="500px" h="300px">
      <CardHeader>{video.title}</CardHeader>

      <CardBody>{video.description}</CardBody>

      <CardFooter>{video.views}</CardFooter>
    </Card>
  );
};

type VideoGridProps = {
  videos: Video[];
};

const VideoGrid = ({ videos }: VideoGridProps) => {
  return (
    <SimpleGrid w="100%" columns={3} spacing={4}>
      {videos.map((video) => (
        <VideoGridItem key={video.id} video={video} />
      ))}
    </SimpleGrid>
  );
};

export default VideoGrid;
