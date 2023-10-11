import AppBar from "@/widgets/app-bar";
import Layout from "@/widgets/layout";
import SideBar from "@/widgets/side-bar";
import { Typography } from "@mui/joy";
import { useParams } from "react-router-dom";

export function Profile() {
  const { userNickname } = useParams();

  return (
    <Layout>
      <AppBar />

      <SideBar />

      <Typography
        level="title-lg"
        textAlign="center"
        style={{
          marginTop: 20,
        }}
      >
        Profile {userNickname}
      </Typography>
    </Layout>
  );
}
