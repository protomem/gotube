import React from "react";
import { Video } from "@/domain/entities";
import { formatDate, formatViews } from "@/lib/utils";

import {
  AspectRatio,
  Card,
  CardFooter,
  SimpleGrid,
  Img,
  Heading,
  HStack,
  Avatar,
  Box,
  Text,
} from "@chakra-ui/react";

type VideoGridItemProps = {
  video: Video;
};

const VideoGridItem = ({ video }: VideoGridItemProps) => {
  return (
    <Card minW="300px" w="auto" maxW="500px">
      <AspectRatio ratio={16 / 9}>
        <Img
          src="https://bit.ly/dan-abramov"
          alt="Dan Abramov"
          roundedTop="md"
        />
      </AspectRatio>

      <CardFooter p={4} pl={6}>
        <HStack>
          <Avatar name="Dan Abramov" w="42px" h="42px" alignSelf="start" />

          <Box>
            <Heading fontSize="lg">{video.title}</Heading>
            <Text fontSize="md">{video.author.nickname}</Text>
            <Text fontSize="sm">
              {`${formatDate(new Date("2023-11-24"))} • ${formatViews(
                video.views,
              )}`}
            </Text>
          </Box>
        </HStack>
      </CardFooter>
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
