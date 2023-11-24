import { videos } from "./fixtures/videos";

export const videoService = {
  getVideo(id: string) {
    return videos[0];
  },

  getVideos() {
    return videos;
  },
};
