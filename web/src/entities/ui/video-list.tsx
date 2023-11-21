import { VideoEntity } from "@/entities/domain/models";
import { VideoCard } from "./video-card";

interface VideoListProps {
  videos: VideoEntity[];
}

export function VideoList({ videos }: VideoListProps) {
  return (
    <div className="flex flex-col gap-5">
      {videos.map((video) => (
        <VideoCard key={video.id} video={video} composit="horizontal" />
      ))}
    </div>
  );
}
