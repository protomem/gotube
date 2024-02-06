import { Video } from "@/domain/entites";
import { SimpleGrid } from "@chakra-ui/react";
import VideoGridItem from "@/components/video-grid-item";

interface Props {
  videos?: Video[];
}

export default function VideoGrid({ videos = [] }: Props) {
  return (
    <SimpleGrid columns={[1, 1, 2, 2, 3, 4]} spacing="6">
      {videos.map((video, index) => (
        <VideoGridItem key={index} video={video} />
      ))}
    </SimpleGrid>
  );
}
