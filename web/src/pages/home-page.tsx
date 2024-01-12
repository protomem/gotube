import { useEffect } from "react";
import { useSearchParams } from "react-router-dom";
import { NavMenuItem, resolveNavMenuItem } from "../components/nav-menu";
import { useQuery } from "@tanstack/react-query";
import { videoService } from "../domain/video.service";
import MainLayout from "../layouts/main-layout";
import VideoGrid from "../components/video-grid";
import { Center, useToast } from "@chakra-ui/react";

const HomePage = () => {
  const toast = useToast();

  const [searchParams] = useSearchParams();
  const nav = resolveNavMenuItem(searchParams.get("nav"));

  const {
    data: videos,
    isError,
    error,
  } = useQuery({
    queryKey: ["videos", nav],
    queryFn: async () => {
      const offset = 0;
      const limit = 10;

      return await videoService.getVideos({
        type: nav === NavMenuItem.Home ? "new" : "popular",
        offset,
        limit,
      });
    },
    select: (data) => data.data.videos,
  });

  useEffect(() => {
    if (isError && error) {
      toast({
        title: "Load Videos Failed",
        description: error.message,
        duration: 3000,
      });
    }
  }, [isError, error]);

  return (
    <MainLayout selectedNavMenuItem={nav}>
      <Center my="6" mx="4">
        <VideoGrid
          videos={videos || []}
          onLastItem={() => console.log("Last Item!!!")}
        />
      </Center>
    </MainLayout>
  );
};

export default HomePage;
