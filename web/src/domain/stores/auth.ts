import { User } from "@/domain/entites";
import { create } from "zustand";
import authStorage from "@/domain/storage/auth.storage";

interface AuthStore {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  loadLocalData: () => void;
  login: (user: User, accessToken: string, refreshToken: string) => void;
  logout: () => void;
  isAuthenticated: () => boolean;
}

export const useAuthStore = create<AuthStore>((set, get) => ({
  user: null,
  accessToken: null,
  refreshToken: null,
  loadLocalData: () => {
    set({
      user: authStorage.getUser(),
      accessToken: authStorage.getAccessToken(),
      refreshToken: authStorage.getRefreshToken(),
    });
  },
  login: (user: User, accessToken: string, refreshToken: string) => {
    set({ user, accessToken, refreshToken });
    authStorage.setUser(user);
    authStorage.setAccessToken(accessToken);
    authStorage.setRefreshToken(refreshToken);
  },
  logout: () => {
    set({ user: null, accessToken: null, refreshToken: null });
    authStorage.removeUser();
    authStorage.removeAccessToken();
    authStorage.removeRefreshToken();
  },
  isAuthenticated: () =>
    get().user !== null &&
    get().accessToken !== null &&
    get().refreshToken !== null,
}));
