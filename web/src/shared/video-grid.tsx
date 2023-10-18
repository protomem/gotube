import { Video } from "@/entities/models";
import { Grid } from "@mui/joy";
import VideoCard from "./video-card";

export interface VideoGridProps {
  videos: Video[];
}

export default function VideoGrid({ videos }: VideoGridProps) {
  return (
    <Grid container spacing={5}>
      {videos.map((video) => (
        <Grid>
          <VideoCard key={video.id} video={video} />
        </Grid>
      ))}
    </Grid>
  );
}
