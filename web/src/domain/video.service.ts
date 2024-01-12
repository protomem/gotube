import { apiClient } from "./api.client";
import { Video } from "./entities";

type GetVideosRequest = {
  limit: number;
  offset: number;
  type: "new" | "popular";
};

type GetVideosResponse = {
  videos: Video[];
};

export const videoService = {
  async getVideos({ limit, offset, type }: GetVideosRequest) {
    return apiClient.get<GetVideosResponse>(
      `/videos/${type}?limit=${limit}&offset=${offset}`,
    );
  },
};
