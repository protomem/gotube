import axios from "axios";
import { authService } from "./auth.service";

export const apiClient = axios.create({
  baseURL: "http://localhost:8080/api", // TODO: from ENV
});

apiClient.interceptors.request.use(
  async (request) => {
    const accessToken = localStorage.getItem("accessToken");

    if (accessToken) {
      request.headers.set("Authorization", `Bearer ${accessToken}`);
    }

    return request;
  },
  (error) => Promise.reject(error),
);

apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const config = error?.config;

    const refreshToken = localStorage.getItem("refreshToken");

    if (error?.response?.status === 401 && !!refreshToken && !config?.sent) {
      config.sent = true;

      const result = await authService.refreshToken({ refreshToken });

      if (result.status === 200) {
        localStorage.setItem("accessToken", result.data.accessToken);
        localStorage.setItem("refreshToken", result.data.refreshToken);

        config.headers["Authorization"] = `Bearer ${result.data.accessToken}`;

        return axios(config);
      }
    }

    return Promise.reject(error);
  },
);

export function resolveAddr(addr: string) {
  if (addr.startsWith("/")) return `${apiClient.defaults.baseURL}${addr}`;
  return addr;
}
