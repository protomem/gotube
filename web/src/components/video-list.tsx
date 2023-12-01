import React, { useRef, useEffect } from "react";
import { Video } from "@/domain/entities";
import { formatDate, formatViews } from "@/lib/utils";
import { ROUTES } from "@/lib/routes";

import NextLink from "next/link";
import {
  AspectRatio,
  Avatar,
  Box,
  Card,
  CardFooter,
  Divider,
  Heading,
  Img,
  Link,
  LinkBox,
  LinkOverlay,
  List,
  Text,
  forwardRef,
} from "@chakra-ui/react";

type VideoListItemProps = {
  video: Video;
};

const VideoListItem = forwardRef(({ video }: VideoListItemProps, ref) => {
  return (
    <LinkBox as={Card} flexDir="row" ref={ref}>
      <AspectRatio width="340px" ratio={16 / 9}>
        <Img src={video.thumbnailPath} alt={video.title} roundedLeft="md" />
      </AspectRatio>

      <Box>
        <Divider orientation="vertical" />
      </Box>

      <CardFooter display="flex" flexDir="column">
        <LinkOverlay as={NextLink} href={`${ROUTES.WATCH}/${video.id}`}>
          <Heading fontSize="lg">{video.title}</Heading>
        </LinkOverlay>
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
          <Link
            as={NextLink}
            href={`${ROUTES.PROFILE}/${video.author.nickname}`}
          >
            <Text fontSize="lg">{video.author.nickname}</Text>
          </Link>
        </Box>
        <Text mt={4}>{video.description}</Text>
      </CardFooter>
    </LinkBox>
  );
});

type VideoListProps = {
  videos: Video[];
  onLast?: () => void;
};

const VideoList = ({ videos, onLast }: VideoListProps) => {
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
    <List spacing={3}>
      {videos.map((video, index) => (
        <VideoListItem
          key={video.id}
          video={video}
          ref={index === videos.length - 1 ? observerTarget : undefined}
        />
      ))}
    </List>
  );
};

export default VideoList;
