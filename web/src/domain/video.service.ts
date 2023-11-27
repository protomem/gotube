import { apiClient } from "@/domain/api.client";
import { withQuery } from "@/lib/routes";
import { Video as VideoEntity } from "@/domain/entities";

export type GetVideosParams = {
  limit?: number;
  offset?: number;
};

type Videos = {
  videos: VideoEntity[];
};

type Video = {
  video: VideoEntity;
};

export const videoService = {
  async getVideo(id: string) {
    return await apiClient.get<Video>(`/videos/${id}`);
  },

  async getNewVideos({ limit, offset }: GetVideosParams) {
    if (limit === undefined) limit = 10;
    if (offset === undefined) offset = 0;

    return await apiClient.get<Videos>(
      withQuery("/videos/new", {
        limit: `${limit}`,
        offset: `${offset}`,
      }),
    );
  },

  async getPopularVideos({ limit, offset }: GetVideosParams) {
    if (limit === undefined) limit = 10;
    if (offset === undefined) offset = 0;

    return await apiClient.get<Videos>(
      withQuery("/videos/popular", {
        limit: `${limit}`,
        offset: `${offset}`,
      }),
    );
  },

  async searchVideos(query: string, { limit, offset }: GetVideosParams) {
    if (limit === undefined) limit = 10;
    if (offset === undefined) offset = 0;

    return await apiClient.get<Videos>(
      withQuery("/videos/search", {
        query,
        limit: `${limit}`,
        offset: `${offset}`,
      }),
    );
  },

  async getUserVideos(nickname: string) {
    return await apiClient.get<Videos>(`/users/${nickname}/videos`);
  },
};
