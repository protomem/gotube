import { apiClient } from "@/entities/api";
import { Video } from "@/entities/models";

export interface GetVideosByUserNicknameRequest {
  userNickname: string;
}

export interface GetVideosByUserNicknameResponse {
  videos: Video[];
}

export const videoService = {
  getVideosByUserNickname: async ({
    userNickname,
  }: GetVideosByUserNicknameRequest) => {
    const response = await apiClient.get<GetVideosByUserNicknameResponse>(
      `/videos/?user=${userNickname}`,
    );
    return response.data;
  },
};
