import { VideoEntity } from "@/entities/domain/models";
import { users } from "./users";

export const videos: VideoEntity[] = [
  {
    id: "1",
    createdAt: new Date(),
    updatedAt: new Date(),
    title: "Title 1",
    description: "Description 1",
    thumbnailPath: "https://img.youtube.com/vi/ZQc9Ez6YpTc/sddefault.jpg",
    videoPath: "https://www.youtube.com/watch?v=ZQc9Ez6YpTc",
    author: users[0],
    isPublic: true,
    views: 1,
  },
];
