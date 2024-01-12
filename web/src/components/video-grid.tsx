import { useEffect, useRef } from "react";
import { Video } from "../domain/entities";
import { Link as RouterLink } from "react-router-dom";
import { formatDate, formatViews } from "../libs";
import { resolveAddr } from "../domain/api.client";
import {
  AspectRatio,
  Card,
  Heading,
  Text,
  Img,
  Link,
  LinkBox,
  LinkOverlay,
  forwardRef,
  CardFooter,
  Divider,
  SimpleGrid,
  HStack,
  Avatar,
  Box,
} from "@chakra-ui/react";

type ItemProps = {
  video: Video;
};

const VideoGridItem = forwardRef(({ video }: ItemProps, ref) => {
  return (
    <LinkBox as={Card} ref={ref} minW="300px" w="auto" maxW="500px">
      <AspectRatio ratio={16 / 9} borderBottom="papayawhip">
        <Img
          src={resolveAddr(video.thumbnailPath)}
          alt={video.title}
          roundedTop="md"
        />
      </AspectRatio>

      <Divider />

      <CardFooter>
        <HStack align="start" gap="2">
          <Link as={RouterLink} to={`profile/${video.author.nickname}`}>
            <Avatar
              src={resolveAddr(video.author.avatarPath)}
              name={video.author.nickname}
              size="md"
              alignSelf="start"
            />
          </Link>

          <Box>
            <LinkOverlay as={RouterLink} to={`watch/${video.id}`}>
              <Heading fontSize="lg">{video.title}</Heading>
            </LinkOverlay>

            <Link as={RouterLink} to={`profile/${video.author.nickname}`}>
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

type Props = {
  videos: Video[];
  onLastItem?: () => void;
};

const VideoGrid = ({ videos, onLastItem }: Props) => {
  const observerTarget = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        const [entry] = entries;
        if (entry.isIntersecting) onLastItem?.();
      },
      { root: null, rootMargin: "0px", threshold: 1 },
    );

    if (observerTarget.current) observer.observe(observerTarget.current);

    return () => {
      if (observerTarget.current) observer.unobserve(observerTarget.current);
    };
  }, [observerTarget, onLastItem]);

  return (
    <SimpleGrid width="full" columns={{ md: 2, lg: 3, "2xl": 4 }} spacing="6">
      {videos.map((video, index) => (
        <VideoGridItem
          key={video.id}
          video={video}
          ref={index === videos.length - 1 ? observerTarget : null}
        />
      ))}
    </SimpleGrid>
  );
};

export default VideoGrid;
