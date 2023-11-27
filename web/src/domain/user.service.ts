import { apiClient } from "@/domain/api.client";
import { User as UserEntity } from "@/domain/entities";

type User = {
  user: UserEntity;
};

export const userService = {
  async getUser(nickname: string) {
    return await apiClient.get<User>(`/users/${nickname}`);
  },
};
