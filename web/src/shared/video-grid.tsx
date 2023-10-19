import { Video } from "@/entities/models";
import { Grid } from "@mui/joy";
import VideoCard from "./video-card";

export interface VideoGridProps {
  videos: Video[];
}

export default function VideoGrid({ videos }: VideoGridProps) {
  return (
    <Grid container spacing={3} xs={12}>
      {videos.map((video) => (
        <Grid sm={4} md={3} key={video.id}>
          <VideoCard key={video.id} video={video} />
        </Grid>
      ))}
    </Grid>
  );
}
