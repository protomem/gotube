import { useRouter } from "next/router";
import { useQuery } from "@tanstack/react-query";
import { videoService } from "@/domain/video.service";
import dynamic from "next/dynamic";

import NextLink from "next/link";
import AppBar from "@/components/app-bar";
import MainLayout from "@/components/layouts/main-layout";
import {
  Avatar,
  Box,
  Button,
  ButtonGroup,
  Card,
  CardBody,
  CardHeader,
  Divider,
  Heading,
  Link,
  LinkBox,
  LinkOverlay,
  Text,
} from "@chakra-ui/react";
import { ROUTES } from "@/lib/routes";

const DynamicVideoPlayer = dynamic(() => import("@/components/video-player"), {
  ssr: false,
});

export default function Watch() {
  const router = useRouter();
  const { videoId } = router.query;

  const { data: video } = useQuery({
    queryKey: ["video", videoId],
    queryFn: async () => {
      return await videoService.getVideo(videoId as string);
    },
    select: (data) => data.data.video,
    enabled: !!videoId,
  });

  return (
    <MainLayout appbar=<AppBar />>
      <Box w="auto" height="full" sx={{ overflowY: "auto" }}>
        <DynamicVideoPlayer
          src={
            "https://storage.googleapis.com/gtv-videos-bucket/sample/ElephantsDream.mp4"
          }
        />

        <Divider my={5} />

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
                <LinkOverlay
                  as={NextLink}
                  href={`${ROUTES.PROFILE}/${"roman"}`}
                >
                  <Heading fontSize="lg">{"Some Author"}</Heading>
                </LinkOverlay>
                <Text>{"324 subscribers"}</Text>
              </Box>

              {/* TODO: hidden on unauthorized */}
              <Button ml={10}>{"Subscribe"}</Button>
            </LinkBox>

            {/* TODO: switch variant   */}
            <Box>
              <ButtonGroup isAttached colorScheme="teal">
                <Button variant="solid">{"2323 likes"}</Button>
                <Button variant="outline">{"323 dislikes"}</Button>
              </ButtonGroup>
            </Box>
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

        <Divider my={5} />

        <Box>{"Comments ..."}</Box>
      </Box>
    </MainLayout>
  );
}
