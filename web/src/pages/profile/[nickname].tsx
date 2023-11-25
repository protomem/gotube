import React from "react";
import { userService } from "@/domain/user.service";
import { videoService } from "@/domain/video.service";
import { repeat } from "@/lib/utils";

import MainLayout from "@/components/layouts/main-layout";
import AppBar from "@/components/app-bar";
import Sidebar from "@/components/side-bar";
import { Box, Divider } from "@chakra-ui/react";
import ProfileCard from "@/components/profile-card";
import VideoGrid from "@/components/video-grid";
import HomeNavMenu from "@/components/home-nav-menu";
import SubscribeButton from "@/components/subscribe-button";

export default function Profile() {
  const [isSubscribed, setIsSubscribed] = React.useState(false);

  const user = userService.getUser("");
  const videos = videoService.getVideos();

  return (
    <MainLayout appbar=<AppBar /> sidebar=<Sidebar navmenu=<HomeNavMenu /> />>
      <Box w="auto" height="full" px={5} py={3} sx={{ overflowY: "auto" }}>
        <ProfileCard
          user={user}
          buttonSubscribe=<SubscribeButton
            isSubscribed={isSubscribed}
            onSubscribe={() => setIsSubscribed(true)}
            onUnsubscribe={() => setIsSubscribed(false)}
          />
        />

        <Divider my={5} />

        <VideoGrid videos={repeat(videos, 13)} />
      </Box>
    </MainLayout>
  );
}
