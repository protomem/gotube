import { User } from "@/domain/entites";
import apiClient from "@/lib/api";

type SignUpRequest = {
  nickname: string;
  email: string;
  password: string;
};

type SignUpResponse = {
  user: User;
  accessToken: string;
  refreshToken: string;
};

type SignInRequest = {
  email: string;
  password: string;
};

type SignInResponse = {
  user: User;
  accessToken: string;
  refreshToken: string;
};

const authService = {
  async signUp(req: SignUpRequest) {
    return await apiClient.post<SignUpResponse>("/auth/register", req);
  },

  async signIn(req: SignInRequest) {
    return await apiClient.post<SignInResponse>("/auth/login", req);
  },
};

export default authService;
