import AppBar from "@/widgets/app-bar";
import Layout from "@/widgets/layout";
import SideBar from "@/widgets/side-bar";
import { Typography } from "@mui/joy";

export default function Home() {
  return (
    <Layout>
      <AppBar />

      <SideBar />

      <Typography
        fontSize={25}
        fontWeight={"bold"}
        textAlign={"center"}
        style={{ marginTop: 20 }}
      >
        Home
      </Typography>
    </Layout>
  );
}
