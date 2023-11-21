import { useSearchParams } from "next/navigation";
import { videoService } from "@/entities/domain/video.service";
import { repeat } from "@/lib";

import { MainLayout } from "@/widgets/layouts/main-layout";
import { AppBar } from "@/widgets/app-bar";
import { SideBar, Navigates } from "@/widgets/side-bar";
import { VideoPane } from "@/widgets/video-pane";

export function HomePage() {
  const searchParams = useSearchParams();

  let selectedNav = Navigates.New;
  if (searchParams && searchParams.has("nav")) {
    switch (searchParams.get("nav")) {
      case Navigates.New:
        selectedNav = Navigates.New;
        break;
      case Navigates.Popular:
        selectedNav = Navigates.Popular;
        break;
      default:
        break;
    }
  }

  const { video } = videoService.getById({ id: "0" });
  const videos = repeat([video], 15);

  return (
    <MainLayout appbar=<AppBar /> sidebar=<SideBar selectedNav={selectedNav} />>
      <div className="w-auto h-full overflow-y-auto">
        <VideoPane videos={videos} composit="grid" />
      </div>
    </MainLayout>
  );
}
