import { Video } from "@/domain/entites";
import apiClient from "@/lib/api";

type FindVideosRequest = {
  limit?: number;
  offset?: number;
};

type FindVideosResponse = {
  videos: Video[];
};

const videoService = {
  async findLatestVideos({ limit = 10, offset = 0 }: FindVideosRequest) {
    return apiClient.get<FindVideosResponse>(
      `/videos?sortBy=latest&limit=${limit}&offset=${offset}`,
    );
  },

  async findTrendingVideos({ limit = 10, offset = 0 }: FindVideosRequest) {
    return apiClient.get<FindVideosResponse>(
      `/videos?sortBy=popular&limit=${limit}&offset=${offset}`,
    );
  },
};

export default videoService;
