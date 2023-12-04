import React from "react";
import { Video } from "@/domain/entities";
import { ROUTES } from "@/lib/routes";

import NextLink from "next/link";
import {
  Avatar,
  Box,
  Heading,
  Text,
  LinkBox,
  LinkOverlay,
  Button,
  ButtonGroup,
  Card,
  CardHeader,
  CardBody,
} from "@chakra-ui/react";

type VideoHeaderProps = {
  video: Video;
  buttonSubscribe?: React.ReactNode;
  buttonRatings: React.ReactNode;
};

const VideoHeader = ({
  video,
  buttonSubscribe,
  buttonRatings,
}: VideoHeaderProps) => {
  return (
    <Box mx={10} display="flex" flexDirection="column" gap={5}>
      <Heading fontSize="2xl">{"Some Title"}</Heading>

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
          <Avatar name="Dan Abrahmov" src="https://bit.ly/dan-abramov" />

          <Box
            display="flex"
            flexDirection="column"
            justifyContent="space-between"
            alignItems="start"
          >
            <LinkOverlay as={NextLink} href={`${ROUTES.PROFILE}/${"roman"}`}>
              <Heading fontSize="lg">{"Some Author"}</Heading>
            </LinkOverlay>
            <Text>{"324 subscribers"}</Text>
          </Box>

          {buttonSubscribe}
        </LinkBox>

        <Box>{buttonRatings}</Box>
      </Box>

      <Card>
        <CardHeader pb={0}>
          <Text fontWeight="bold">{"3232 views • 2 years"}</Text>
        </CardHeader>
        <CardBody pt={0}>
          <Text>{"Some description ..."}</Text>
        </CardBody>
      </Card>
    </Box>
  );
};

export default VideoHeader;
