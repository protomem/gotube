import { AspectRatio, Box, Center } from "@chakra-ui/react";
import ReactVideoPlayer from "react-player";

type VideoPlayerProps = {
  src: string;
};

const VideoPlayer = ({ src }: VideoPlayerProps) => {
  return (
    <Center w="full" bg="black">
      <ReactVideoPlayer url={src} controls width="auto" height="540px" />
    </Center>
  );
};

export default VideoPlayer;
