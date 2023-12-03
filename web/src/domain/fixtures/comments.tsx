import { Comment } from "../entities";
import { users } from "./users";

export const comments: Comment[] = [
  {
    id: "123454434",
    createdAt: new Date(),
    updatedAt: new Date(),
    content: "Some Comment",
    author: users[0],
    videoId: "video1234",
  },
];
