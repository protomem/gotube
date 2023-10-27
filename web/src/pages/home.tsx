import { useState } from "react";
import { NavItem } from "@/feature/nav-menu";
import VideoGrid from "@/feature/video-grid";
import Paginator from "@/shared/paginator";
import AppBar from "@/widgets/app-bar";
import Layout from "@/widgets/layout";
import SideBar from "@/widgets/side-bar";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { videoService } from "@/entities/video.service";
import { Video } from "@/entities/models";
import { useAppSelector } from "@/feature/store/hooks";
import { selectAccessToken } from "@/feature/store/auth/auth.selectors";
import { Box } from "@mui/joy";

export default function Home() {
  const nav = useNavigate();

  const accessToken = useAppSelector(selectAccessToken);
  const [searchParams] = useSearchParams();

  let selectedNavItem: NavItem;
  switch (searchParams.get("videos")) {
    case NavItem.New.toLowerCase():
      selectedNavItem = NavItem.New;
      break;
    case NavItem.Popular.toLowerCase():
      selectedNavItem = NavItem.Popular;
      break;
    case NavItem.Subscriptions.toLowerCase():
      selectedNavItem = NavItem.Subscriptions;
      break;
    default:
      selectedNavItem = NavItem.New;
      break;
  }

  const [maxPage, setMaxPage] = useState(0);
  const pageSize = 16;
  let page = 1;
  if (searchParams.has("page") && Number(searchParams.get("page")) > 0) {
    page = Number(searchParams.get("page"));
  }

  const [videos, setVideos] = useState<Video[]>([]);

  useQuery({
    queryKey: ["videos", selectedNavItem, page],
    queryFn: () => {
      switch (selectedNavItem) {
        case NavItem.New:
          return videoService.getNewVideos({ page, pageSize });
        case NavItem.Popular:
          return videoService.getPopularVideos({ page, pageSize });
        case NavItem.Subscriptions:
          return videoService.getSubscriptionsVideos({
            page,
            pageSize,
            accessToken: accessToken || "",
          });
        default:
          return videoService.getNewVideos({ page, pageSize });
      }
    },
    onSuccess: (data) => {
      setMaxPage(Math.ceil(data.totalCount / pageSize));
      setVideos(data.videos);
    },
    onError: (error) => {
      console.error(error);
      setVideos([]);
    },
    refetchOnMount: true,
  });

  return (
    <Layout>
      <AppBar />

      <SideBar selectedNavItem={selectedNavItem} />

      <Box
        style={{
          height: "90vh",
          display: "flex",
          flexDirection: "column",
          justifyContent: "start",
        }}
      >
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

        <Box style={{ flex: 1, alignSelf: "center", marginTop: 10 }}>
          <Paginator
            page={page}
            onChangePage={(newPage) => {
              nav(`/?page=${newPage}&videos=${selectedNavItem.toLowerCase()}`, {
                replace: true,
              });
            }}
            maxPage={maxPage}
          />
        </Box>
      </Box>
    </Layout>
  );
}
