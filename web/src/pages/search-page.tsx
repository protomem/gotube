import { useSearchParams } from "next/navigation";
import { AppBar } from "@/widgets/app-bar";
import { MainLayout } from "@/widgets/layouts/main-layout";
import { SideBar } from "@/widgets/side-bar";
import { VideoCard } from "@/entities/ui/video-card";
import { videos } from "@/fixtures/videos";

export function SearchPage() {
  let query = "";

  const searchParams = useSearchParams();
  if (searchParams !== null && searchParams.has("q")) {
    query = searchParams.get("q") || "";
  }

  return (
    <MainLayout appbar=<AppBar /> sidebar=<SideBar />>
      <div className="w-auto h-full overflow-y-auto">
        <h1>Search: {query}</h1>

        <div className="flex flex-col gap-5 ml-3 mr-5">
          <VideoCard video={videos[0]} composit="horizontal" />
          <VideoCard video={videos[0]} composit="horizontal" />
          <VideoCard video={videos[0]} composit="horizontal" />
          <VideoCard video={videos[0]} composit="horizontal" />
          <VideoCard video={videos[0]} composit="horizontal" />
          <VideoCard video={videos[0]} composit="horizontal" />
        </div>
      </div>
    </MainLayout>
  );
}
