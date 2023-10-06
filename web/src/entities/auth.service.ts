import { apiClient } from "./api";
import { User } from "./models";

export interface RegisterRequest {
  nickname: string;
  email: string;
  password: string;
}

export interface RegisterResponse {
  accessToken: string;
  refreshToken: string;
  user: User;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  accessToken: string;
  refreshToken: string;
  user: User;
}

export interface LogoutRequest {
  accessToken: string;
  refreshToken: string;
}

export const authService = {
  async register(request: RegisterRequest) {
    const response = await apiClient.post<RegisterResponse>(
      "/auth/register",
      request,
    );
    return response.data;
  },

  async login(request: LoginRequest) {
    const response = await apiClient.post<LoginResponse>(
      "/auth/login",
      request,
    );
    return response.data;
  },

  async logout(request: LogoutRequest) {
    await apiClient.delete(`/auth/logout?session=${request.refreshToken}`, {
      headers: {
        Authorization: `Bearer ${request.accessToken}`,
      },
    });
  },
};
