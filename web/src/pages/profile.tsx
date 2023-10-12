import { useState } from "react";
import { User } from "@/entities/models";
import { userService } from "@/entities/user.service";
import AppBar from "@/widgets/app-bar";
import Layout from "@/widgets/layout";
import SideBar from "@/widgets/side-bar";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import ProfilePane from "@/widgets/profile-pane";
import { Box, Divider } from "@mui/joy";

export function Profile() {
  const { userNickname } = useParams();

  const [user, setUser] = useState<User | null>(null);
  useQuery({
    queryKey: ["users", userNickname],
    queryFn: async () => userService.getUser({ nickname: userNickname || "" }),
    onSuccess: (data) => {
      setUser(data.user);
    },
    enabled: userNickname !== undefined,
  });

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
        </>
      )}
    </Layout>
  );
}
