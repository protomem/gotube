import { apiClient } from "@/entities/api";
import { Video } from "@/entities/models";

export interface GetVideosRequest {
  page: number;
  pageSize: number;
}

export interface GetVideosResponse {
  totalCount: number;
  videos: Video[];
}

export interface GetSubscriptionsVideosRequest extends GetVideosRequest {
  accessToken: string;
}

export interface GetVideosByUserNicknameRequest {
  userNickname: string;
}

export interface GetVideosByUserNicknameResponse {
  videos: Video[];
}

export const videoService = {
  async getNewVideos(request: GetVideosRequest) {
    const response = await apiClient.get<GetVideosResponse>(
      `/videos/?sort_by=new&limit=${request.pageSize}&offset=${
        (request.page - 1) * request.pageSize
      }`,
    );
    return response.data;
  },

  async getPopularVideos(request: GetVideosRequest) {
    const response = await apiClient.get<GetVideosResponse>(
      `/videos/?sort_by=popular&limit=${request.pageSize}&offset=${
        (request.page - 1) * request.pageSize
      }`,
    );
    return response.data;
  },

  async getSubscriptionsVideos(request: GetSubscriptionsVideosRequest) {
    const response = await apiClient.get<GetVideosResponse>(
      `/videos/?is_subs=true&limit=${request.pageSize}&offset=${
        (request.page - 1) * request.pageSize
      }`,
      {
        headers: {
          Authorization: `Bearer ${request.accessToken}`,
        },
      },
    );
    return response.data;
  },

  async getVideosByUserNickname(request: GetVideosByUserNicknameRequest) {
    const response = await apiClient.get<GetVideosByUserNicknameResponse>(
      `/videos/?user=${request.userNickname}`,
    );
    return response.data;
  },
};
