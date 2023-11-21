import React from "react";

import { VideoEntity } from "@/entities/domain/models";
import { chunk } from "@/lib";
import { VideoCard } from "./video-card";

interface VideoGridRowProps {
  videos: VideoEntity[];
}

const VideoGridRow: React.FC<VideoGridRowProps> = ({ videos }) => {
  return (
    <div className="flex flex-row gap-7 justify-start">
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
    <div className="flex flex-col gap-7 justify-center items-start">
      {videosGrid.map((videosRow, index) => (
        <VideoGridRow key={index} videos={videosRow} />
      ))}
    </div>
  );
}
