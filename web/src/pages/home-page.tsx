import MainLayout from "../layouts/main-layout";
import VideoGrid from "../components/video-grid";
import { NavMenuItem } from "../components/nav-menu";

const HomePage = () => {
  return (
    <MainLayout hideSideBar selectedNavMenuItem={NavMenuItem.Home}>
      <VideoGrid />
    </MainLayout>
  );
};

export default HomePage;
