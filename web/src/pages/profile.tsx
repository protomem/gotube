import { useState } from "react";
import { User, Video } from "@/entities/models";
import { userService } from "@/entities/user.service";
import AppBar from "@/widgets/app-bar";
import Layout from "@/widgets/layout";
import SideBar from "@/widgets/side-bar";
import { useQuery } from "@tanstack/react-query";
import { useNavigate, useParams } from "react-router-dom";
import ProfilePane from "@/widgets/profile-pane";
import { Box, Divider } from "@mui/joy";
import VideoGrid from "@/feature/video-grid";
import { videoService } from "@/entities/video.service";

export function Profile() {
  const nav = useNavigate();
  const { userNickname } = useParams();

  const [user, setUser] = useState<User | null>(null);
  useQuery({
    queryKey: ["users", userNickname],
    queryFn: async () => userService.getUser({ nickname: userNickname || "" }),
    onSuccess: (data) => {
      setUser(data.user);
    },
    onError: (error) => {
      console.error(error);
      nav("/not-found", { replace: true });
    },
    enabled: userNickname !== undefined,
  });

  const [videos, setVideos] = useState<Video[]>([]);
  useQuery({
    queryKey: ["videos", userNickname],
    queryFn: async () =>
      videoService.getVideosByUserNickname({
        userNickname: user?.nickname || "",
      }),
    onSuccess: (data) => {
      setVideos(data.videos);
    },
    enabled: user !== null,
  });

  return (
    <Layout>
      <AppBar />

      <SideBar />

      {user !== null && (
        <Box
          style={{
            height: "90vh",
            display: "flex",
            flexDirection: "column",
            justifyContent: "start",
          }}
        >
          <Box style={{ flex: 5, padding: 20 }}>
            <ProfilePane user={user} />
          </Box>

          <Divider
            style={{
              marginLeft: 20,
              marginRight: 20,
              height: 4,
            }}
          />

          <Box
            style={{
              flex: 24,
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              marginTop: 30,
              marginLeft: 20,
              marginRight: 20,
              overflowY: "auto",
            }}
          >
            <VideoGrid videos={videos} />
          </Box>
        </Box>
      )}
    </Layout>
  );
}
