import { useSearchParams } from "next/navigation";
import { videoService } from "@/entities/domain/video.service";

import { MainLayout } from "@/widgets/layouts/main-layout";
import { AppBar } from "@/widgets/app-bar";
import { SideBar, Navigates } from "@/widgets/side-bar";
import { VideoCard } from "@/entities/ui/video-card";

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

  return (
    <MainLayout appbar=<AppBar /> sidebar=<SideBar selectedNav={selectedNav} />>
      <div className="w-auto h-full overflow-y-auto">
        <div className="m-5 mt-10 flex justify-around gap-7">
          <VideoCard video={video} />
          <VideoCard video={video} />
          <VideoCard video={video} />
        </div>

        <div className="m-5 mt-7 flex justify-around gap-7">
          <VideoCard video={video} />
          <VideoCard video={video} />
          <VideoCard video={video} />
        </div>

        <div className="m-5 mt-7 flex justify-around gap-7">
          <VideoCard video={video} />
          <VideoCard video={video} />
          <VideoCard video={video} />
        </div>
      </div>
    </MainLayout>
  );
}
