import React from "react";
import { Video } from "@/domain/entities";
import { ROUTES } from "@/lib/routes";
import { formatDate, formatViews } from "@/lib/utils";

import NextLink from "next/link";
import {
  Avatar,
  Box,
  Heading,
  Text,
  LinkBox,
  LinkOverlay,
  Card,
  CardHeader,
  CardBody,
} from "@chakra-ui/react";

type VideoHeaderProps = {
  video: Video;
  subscribers: number;
  buttonSubscribe?: React.ReactNode;
  buttonRatings: React.ReactNode;
};

const VideoHeader = ({
  video,
  subscribers,
  buttonSubscribe,
  buttonRatings,
}: VideoHeaderProps) => {
  return (
    <Box mx={10} display="flex" flexDirection="column" gap={5}>
      <Heading fontSize="2xl">{video.title}</Heading>

      <Box
        display="flex"
        flexDirection="row"
        justifyContent="space-between"
        alignItems="center"
      >
        <LinkBox
          display="flex"
          flexDirection="row"
          justifyContent="space-between"
          alignItems="center"
          gap={3}
        >
          <Avatar name={video.author.nickname} src={video.author.avatarPath} />

          <Box
            display="flex"
            flexDirection="column"
            justifyContent="space-between"
            alignItems="start"
          >
            <LinkOverlay
              as={NextLink}
              href={`${ROUTES.PROFILE}/${video.author.nickname}`}
            >
              <Heading fontSize="lg">{video.author.nickname}</Heading>
            </LinkOverlay>
            <Text>{`${formatViews(subscribers)} subscribers`}</Text>
          </Box>

          {buttonSubscribe}
        </LinkBox>

        <Box>{buttonRatings}</Box>
      </Box>

      <Card>
        <CardHeader pb={0}>
          <Text fontWeight="bold">{`${formatViews(
            video.views,
          )} views • ${formatDate(video.createdAt)}`}</Text>
        </CardHeader>
        <CardBody pt={0}>
          <Text>{video.description}</Text>
        </CardBody>
      </Card>
    </Box>
  );
};

export default VideoHeader;
