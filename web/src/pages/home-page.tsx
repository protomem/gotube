import MainLayout from "../layouts/main-layout";
import VideoGrid from "../components/video-grid";
import { resolveNavMenuItem } from "../components/nav-menu";
import { useSearchParams } from "react-router-dom";
import { Center, Heading } from "@chakra-ui/react";

const HomePage = () => {
  const [searchParams] = useSearchParams();
  const nav = resolveNavMenuItem(searchParams.get("nav"));

  return (
    <MainLayout hideSideBar selectedNavMenuItem={nav}>
      <Center>
        <Heading>{nav} Page</Heading>
      </Center>

      <VideoGrid />
    </MainLayout>
  );
};

export default HomePage;
