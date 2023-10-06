import { RootState } from "@/feature/store/store";

const selectAuthModule = (state: RootState) => state.auth;

export const selectUser = (state: RootState) => selectAuthModule(state).user;
export const selectAccessToken = (state: RootState) =>
  selectAuthModule(state).accessToken;
export const selectRefreshToken = (state: RootState) =>
  selectAuthModule(state).refreshToken;
export const selectIsLoggedIn = (state: RootState) =>
  selectUser(state) !== null &&
  selectAccessToken(state) !== null &&
  selectRefreshToken(state) !== null;
