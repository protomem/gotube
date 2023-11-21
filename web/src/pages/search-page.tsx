import { useSearchParams } from "next/navigation";
import { videos } from "@/fixtures/videos";
import { repeat } from "@/lib";

import { AppBar } from "@/widgets/app-bar";
import { MainLayout } from "@/widgets/layouts/main-layout";
import { SideBar } from "@/widgets/side-bar";
import { VideoPane } from "@/widgets/video-pane";

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

        <div className="ml-3 mr-5">
          <VideoPane videos={repeat(videos, 2)} composit="list" />
        </div>
      </div>
    </MainLayout>
  );
}
