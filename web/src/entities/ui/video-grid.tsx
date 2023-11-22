import React from "react";
import { chunk } from "@/lib";
import { VideoEntity } from "@/entities/domain/models";

import { VideoCard } from "@/entities/ui/video-card";

interface VideoGridRowProps {
  videos: VideoEntity[];
}

const VideoGridRow: React.FC<VideoGridRowProps> = ({ videos }) => {
  return (
    <div className="flex flex-row gap-7 justify-start flex-3">
      {videos.map((video) => (
        <VideoCard key={video.id} video={video} composit="vertical" />
      ))}
    </div>
  );
};

interface VideoGridProps {
  videos: VideoEntity[];
  column?: number;
}

export function VideoGrid({ videos, column }: VideoGridProps) {
  if (column === undefined) column = 3;
  const videosGrid = chunk(videos, column);

  return (
    <div className="w-full flex flex-col gap-7 justify-center">
      {videosGrid.map((videosRow, index) => (
        <VideoGridRow key={index} videos={videosRow} />
      ))}
    </div>
  );
}
