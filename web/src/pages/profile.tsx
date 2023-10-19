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
import VideoGrid from "@/shared/video-grid";

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

  const videos: Video[] = [
    {
      id: "1",
      createdAt: new Date(),
      updatedAt: new Date(),
      title: "Video 1",
      description: "Description 1",
      thumbnailPath:
        "https://images.unsplash.com/photo-1532614338840-ab30cf10ed36?auto=format&fit=crop&w=318",
      videoPath: "https://example.com/video1.mp4",
      isPublic: true,
      views: 100,
      user: {
        id: "1",
        createdAt: new Date(),
        updatedAt: new Date(),
        nickname: "user1",
        email: "0pEeH@example.com",
        isVerified: true,
        avatarPath: "https://example.com/user1.jpg",
        description: "Description 1",
      },
    },
    {
      id: "2",
      createdAt: new Date(),
      updatedAt: new Date(),
      title: "Video 2",
      description: "Description 2",
      thumbnailPath:
        "https://images.unsplash.com/photo-1532614338840-ab30cf10ed36?auto=format&fit=crop&w=318",
      videoPath: "https://example.com/video2.mp4",
      isPublic: true,
      views: 200,
      user: {
        id: "2",
        createdAt: new Date(),
        updatedAt: new Date(),
        nickname: "user2",
        email: "0pEeH@example.com",
        isVerified: true,
        avatarPath: "https://example.com/user2.jpg",
        description: "Description 2",
      },
    },
    {
      id: "3",
      createdAt: new Date(),
      updatedAt: new Date(),
      title: "Video 3",
      description: "Description 3",
      thumbnailPath:
        "https://images.unsplash.com/photo-1532614338840-ab30cf10ed36?auto=format&fit=crop&w=318",
      videoPath: "https://example.com/video3.mp4",
      isPublic: true,
      views: 300,
      user: {
        id: "3",
        createdAt: new Date(),
        updatedAt: new Date(),
        nickname: "user3",
        email: "0pEeH@example.com",
        isVerified: true,
        avatarPath: "https://example.com/user3.jpg",
        description: "Description 3",
      },
    },
    {
      id: "4",
      createdAt: new Date(),
      updatedAt: new Date(),
      title: "Video 4",
      description: "Description 4",
      thumbnailPath:
        "https://images.unsplash.com/photo-1532614338840-ab30cf10ed36?auto=format&fit=crop&w=318",
      videoPath: "https://example.com/video4.mp4",
      isPublic: true,
      views: 400,
      user: {
        id: "4",
        createdAt: new Date(),
        updatedAt: new Date(),
        nickname: "user4_ddddd_daaaaaaaaaaaa",
        email: "0pEeH@example.com",
        isVerified: true,
        avatarPath: "https://example.com/user4.jpg",
        description: "Description 4",
      },
    },
  ];

  return (
    <Layout>
      <AppBar />

      <SideBar />

      {user !== null && (
        <>
          <Box style={{ padding: 20 }}>
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
              display: "flex",
              flexDirection: "column",
              justifyContent: "center",
              alignItems: "center",
              marginTop: 30,
              marginLeft: 20,
              marginRight: 20,
            }}
          >
            <VideoGrid videos={videos} />
          </Box>
        </>
      )}
    </Layout>
  );
}
