import { useSearchParams } from "next/navigation";
import { videos } from "@/fixtures/videos";
import { repeat } from "@/lib";

import { AppBar } from "@/widgets/app-bar";
import { MainLayout } from "@/widgets/layouts/main-layout";
import { SideBar } from "@/widgets/side-bar";
import { VideoPane } from "@/widgets/video-pane";
import { Separator } from "@/shared/ui/separator";

export function SearchPage() {
  let query = "";

  const searchParams = useSearchParams();
  if (searchParams !== null && searchParams.has("q")) {
    query = searchParams.get("q") || "";
  }

  return (
    <MainLayout appbar=<AppBar /> sidebar=<SideBar />>
      <div className="w-auto h-full overflow-y-auto">
        <div className="m-8 ml-3 mr-5">
          <div className="ml-2 flex flex-row gap-3 justify-start items-baseline">
            <h1 className="text-xl font-bold">Search: </h1>
            <h3 className="text-lg italic">{query}</h3>
          </div>

          <Separator className="mt-3 mb-5" />

          <VideoPane videos={repeat(videos, 20)} composit="list" />
        </div>
      </div>
    </MainLayout>
  );
}
