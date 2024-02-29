import { Video } from "@/domain/entites";
import NextLink from "next/link";
import {
  AspectRatio,
  Avatar,
  Card,
  CardBody,
  HStack,
  Heading,
  Text,
  Img,
  VStack,
  LinkBox,
  LinkOverlay,
  Link,
} from "@chakra-ui/react";

interface Props {
  video: Video;
}

export default function VideoGridItem({ video }: Props) {
  return (
    <Card minW="xs" w="100%" variant="ghost">
      <CardBody as={LinkBox} p="0" pb="2">
        <AspectRatio ratio={16 / 9}>
          <Img
            rounded="md"
            src={`${process.env.NEXT_PUBLIC_API_ADDR}/${video.thumbnailPath}`}
          />
        </AspectRatio>

        <HStack mt="4" mx="2" alignItems="start">
          <Link as={NextLink} href={`/profile/${video.author.nickname}`}>
            <Avatar
              name={video.author.nickname}
              src={`${process.env.NEXT_PUBLIC_API_ADDR}/${video.author.avatarPath}`}
            />
          </Link>

          <VStack alignItems="start">
            <LinkOverlay as={NextLink} href={`/video/${video.id}`}>
              <Heading size="md">{video.title}</Heading>
            </LinkOverlay>

            <Text noOfLines={[1, 2]}>{video.description || "..."}</Text>

            <Text>
              {`${video.views} views • ${new Date(video.createdAt).toLocaleTimeString().split(":").slice(0, 2).join(":")} ${new Date(video.createdAt).toLocaleTimeString().split(" ")[1]} ${new Date(video.createdAt).toLocaleDateString()}`}
            </Text>
          </VStack>
        </HStack>
      </CardBody>
    </Card>
  );
}
