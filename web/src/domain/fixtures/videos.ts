import { Video } from "../entities";
import { users } from "./users";

export const videos: Video[] = [
  {
    id: "1",
    createdAt: new Date(),
    updatedAt: new Date(),
    title: "Video 1",
    description: "Description 1",
    thumbnailPath: "https://example.com/thumbnail1.jpg",
    videoPath: "https://example.com/video1.mp4",
    author: users[0],
    isPublic: true,
    views: 23433,
  },
];
