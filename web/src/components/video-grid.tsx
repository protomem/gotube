import React, { useEffect, useRef } from "react";
import { Video } from "@/domain/entities";
import { formatDate, formatViews } from "@/lib/utils";
import { ROUTES } from "@/lib/routes";

import NextLink from "next/link";
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
  Divider,
  LinkOverlay,
  Link,
  LinkBox,
  forwardRef,
} from "@chakra-ui/react";

type VideoGridItemProps = {
  video: Video;
};

const VideoGridItem = forwardRef(({ video }: VideoGridItemProps, ref) => {
  return (
    <LinkBox as={Card} minW="300px" w="auto" maxW="500px" ref={ref}>
      <AspectRatio ratio={16 / 9} borderBottom="papayawhip">
        <Img src={video.thumbnailPath} alt={video.title} roundedTop="md" />
      </AspectRatio>

      <Divider />

      <CardFooter p={4} pl={6}>
        <HStack>
          <Avatar
            src={video.author.avatarPath}
            name={video.author.nickname}
            w="42px"
            h="42px"
            alignSelf="start"
          />

          <Box>
            <LinkOverlay as={NextLink} href={`${ROUTES.WATCH}/${video.id}`}>
              <Heading fontSize="lg">{video.title}</Heading>
            </LinkOverlay>
            <Link
              as={NextLink}
              href={`${ROUTES.PROFILE}/${video.author.nickname}`}
            >
              <Text fontSize="md">{video.author.nickname}</Text>
            </Link>
            <Text fontSize="sm">
              {`${formatDate(video.createdAt)} • ${formatViews(
                video.views,
              )} views`}
            </Text>
          </Box>
        </HStack>
      </CardFooter>
    </LinkBox>
  );
});

type VideoGridProps = {
  videos: Video[];
  onLast?: () => void;
};

const VideoGrid = ({ videos, onLast }: VideoGridProps) => {
  const observerTarget = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        const [entry] = entries;
        if (entry.isIntersecting) onLast?.();
      },
      { root: null, rootMargin: "0px", threshold: 1 },
    );

    if (observerTarget.current) observer.observe(observerTarget.current);

    return () => {
      if (observerTarget.current) observer.unobserve(observerTarget.current);
    };
  }, [observerTarget, onLast]);

  return (
    <SimpleGrid w="100%" columns={3} spacing={4}>
      {videos.map((video, index) => (
        <VideoGridItem
          key={video.id}
          video={video}
          ref={index === videos.length - 1 ? observerTarget : undefined}
        />
      ))}
    </SimpleGrid>
  );
};

export default VideoGrid;
