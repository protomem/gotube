import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "next/router";
import { userService } from "@/domain/user.service";
import { videoService } from "@/domain/video.service";

import MainLayout from "@/components/layouts/main-layout";
import AppBar from "@/components/app-bar";
import Sidebar from "@/components/side-bar";
import { Box, Divider } from "@chakra-ui/react";
import ProfileCard from "@/components/profile-card";
import VideoGrid from "@/components/video-grid";
import HomeNavMenu from "@/components/home-nav-menu";
import SubscribeButton from "@/components/subscribe-button";
import { User } from "@/domain/entities";

export default function Profile() {
  const router = useRouter();
  const { nickname } = router.query;

  const [isSubscribed, setIsSubscribed] = useState(false);

  const { data: user } = useQuery({
    queryKey: ["user", nickname],
    queryFn: () => userService.getUser(nickname as string),
    enabled: !!nickname,
    select: (data) => data.data.user,
  });

  const { data: videos } = useQuery({
    queryKey: ["videos", nickname],
    queryFn: () => videoService.getUserVideos(nickname as string),
    enabled: !!nickname,
    select: (data) => data.data.videos,
  });

  return (
    <MainLayout appbar=<AppBar /> sidebar=<Sidebar navmenu=<HomeNavMenu /> />>
      <Box w="auto" height="full" px={5} py={3} sx={{ overflowY: "auto" }}>
        <ProfileCard user={user ?? ({} as User)} />

        <Divider my={5} />

        <VideoGrid videos={videos ?? []} />
      </Box>
    </MainLayout>
  );
}
