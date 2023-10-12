import { apiClient } from "@/entities/api";
import { User } from "@/entities/models";

export interface GetUserRequest {
  nickname: string;
}

export interface GetUserResponse {
  user: User;
}

export const userService = {
  async getUser(request: GetUserRequest) {
    const response = await apiClient.get<GetUserResponse>(
      `/users/${request.nickname}`,
    );
    return response.data;
  },
};
