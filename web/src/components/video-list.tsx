import { Video } from "@/domain/entites";
import { List, ListItem } from "@chakra-ui/react";
import VideoListItem from "./video-list-item";

interface Props {
  videos?: Video[];
}

export default function VideoList({ videos = [] }: Props) {
  return (
    <List spacing={6}>
      {videos.map((video) => (
        <ListItem key={video.id}>
          <VideoListItem video={video} />
        </ListItem>
      ))}
    </List>
  );
}
