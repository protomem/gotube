import { User } from "@/domain/entites";
import { create } from "zustand";

interface AuthStore {
  user?: User;
  accessToken?: string;
  refreshToken?: string;
}

// TODO: add persisted state
export const useAuthStore = create<AuthStore>((set, get) => ({
  user: undefined,
  accessToken: undefined,
  refreshToken: undefined,
  login: (user: User, accessToken: string, refreshToken: string) =>
    set({ user, accessToken, refreshToken }),
  logout: () =>
    set({ user: undefined, accessToken: undefined, refreshToken: undefined }),
  isAuthenticated: () =>
    get().user !== undefined &&
    get().accessToken !== undefined &&
    get().refreshToken !== undefined,
}));
