import { apiClient } from "@/domain/api.client";
import { Comment as CommentEntity } from "@/domain/entities";

type Comments = {
  comments: CommentEntity[];
};

export const commentService = {
  async getCommentsByVideo(videoId: string) {
    return await apiClient.get<Comments>(`/videos/${videoId}/comments`);
  },
};
