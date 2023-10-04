import { RootState } from "@/feature/store/store";

const selectAuthModule = (state: RootState) => state.auth;

export const selectUser = (state: RootState) => selectAuthModule(state).user;
export const selectAccessToken = (state: RootState) =>
  selectAuthModule(state).accessToken;
