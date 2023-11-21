import { VideoEntity } from "@/entities/domain/models";
import { VideoGrid } from "@/entities/ui/video-grid";
import { VideoList } from "@/entities/ui/video-list";

interface VideoPaneProps {
  videos?: VideoEntity[];
  composit?: "grid" | "list";
}

export function VideoPane({ videos, composit }: VideoPaneProps) {
  if (videos === undefined) videos = [];
  if (composit === undefined) composit = "grid";

  return (
    <div>
      {composit === "grid" ? (
        <VideoGrid videos={videos} />
      ) : (
        <VideoList videos={videos} />
      )}
    </div>
  );
}
