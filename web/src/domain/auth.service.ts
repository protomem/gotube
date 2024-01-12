import { apiClient } from "./api.client";
import { User } from "./entities";

type LoginRequest = {
  email: string;
  password: string;
};

type LoginResponse = {
  accessToken: string;
  refreshToken: string;
  user: User;
};

type RegisterRequest = {
  nickname: string;
  email: string;
  password: string;
};

type RegisterResponse = {
  accessToken: string;
  refreshToken: string;
  user: User;
};

type RefreshTokenRequest = {
  refreshToken: string;
};

type RefreshTokenResponse = {
  accessToken: string;
  refreshToken: string;
};

export const authService = {
  async login({ email, password }: LoginRequest) {
    return await apiClient.post<LoginResponse>("/auth/login", {
      email,
      password,
    });
  },

  async register({ nickname, email, password }: RegisterRequest) {
    return await apiClient.post<RegisterResponse>("/auth/register", {
      nickname,
      email,
      password,
    });
  },

  async refreshToken({ refreshToken }: RefreshTokenRequest) {
    return await apiClient.get<RefreshTokenResponse>("/auth/refresh", {
      headers: {
        "X-Refresh-Token": refreshToken,
      },
    });
  },
};
