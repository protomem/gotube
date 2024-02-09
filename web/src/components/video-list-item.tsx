import { Video } from "@/domain/entites";
import NextLink from "next/link";
import {
  AspectRatio,
  Box,
  Card,
  CardBody,
  HStack,
  Heading,
  Text,
  Img,
  LinkOverlay,
  Link,
  Avatar,
  LinkBox,
} from "@chakra-ui/react";

interface Props {
  video: Video;
}

export default function VideoListItem({ video }: Props) {
  return (
    <Card as={LinkBox} direction="row" variant="ghost">
      <AspectRatio ratio={16 / 9} w="sm">
        <Img
          maxW="sm"
          objectFit="cover"
          rounded="md"
          src={`${process.env.NEXT_PUBLIC_API_ADDR}/${video.thumbnailPath}`}
        />
      </AspectRatio>

      <CardBody py="0">
        <LinkOverlay as={NextLink} href={`/video/${video.id}`}>
          <Heading size="lg" fontWeight="normal">
            {video.title}
          </Heading>
        </LinkOverlay>

        <Text>
          {`${video.views} views • ${new Date(video.createdAt).toLocaleTimeString().split(":").slice(0, 2).join(":")} ${new Date(video.createdAt).toLocaleTimeString().split(" ")[1]} ${new Date(video.createdAt).toLocaleDateString()}`}
        </Text>

        <Box display="flex" alignItems="center" gap="2" mt="2">
          <Link as={NextLink} href={`/profile/${video.author.nickname}`}>
            <Avatar
              name={video.author.nickname}
              src={`${process.env.NEXT_PUBLIC_API_ADDR}/${video.author.avatarPath}`}
            />
          </Link>

          <Text fontSize="lg">{video.author.nickname}</Text>
        </Box>

        <Text mt="4" noOfLines={[2]}>
          {video.description || "..."}
        </Text>
      </CardBody>
    </Card>
  );
}
