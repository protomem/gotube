import { User } from "./models";

export const authRepo = {
  setUser(user: User) {
    localStorage.setItem("auth_user", JSON.stringify(user));
  },

  getUser() {
    const userJSON = localStorage.getItem("auth_user");
    if (userJSON === null) return null;
    return JSON.parse(userJSON);
  },

  removeUser() {
    localStorage.removeItem("auth_user");
  },

  setAccessToken(accessToken: string) {
    localStorage.setItem("auth_accessToken", accessToken);
  },

  getAccessToken() {
    return localStorage.getItem("auth_accessToken");
  },

  removeAccessToken() {
    localStorage.removeItem("auth_accessToken");
  },
};
