import { User } from "@/domain/entites";

const authStorage = {
  setAccessToken(accessToken: string) {
    localStorage.setItem("accessToken", accessToken);
  },

  getAccessToken(): string | null {
    return localStorage.getItem("accessToken");
  },

  removeAccessToken() {
    localStorage.removeItem("accessToken");
  },

  setRefreshToken(refreshToken: string) {
    localStorage.setItem("refreshToken", refreshToken);
  },

  getRefreshToken(): string | null {
    return localStorage.getItem("refreshToken");
  },

  removeRefreshToken() {
    localStorage.removeItem("refreshToken");
  },

  setUser(user: User) {
    const userJSON = JSON.stringify(user);
    localStorage.setItem("user", userJSON);
  },

  getUser(): User | null {
    const userJSON = localStorage.getItem("user");
    if (userJSON === null) return null;
    return JSON.parse(userJSON) as User;
  },

  removeUser() {
    localStorage.removeItem("user");
  },
};

export default authStorage;
