import { PayloadAction, createSlice } from "@reduxjs/toolkit";
import { User } from "@/entities/models";
import { authRepo } from "@/entities/auth.repository";

export interface AuthState {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
}

const initialState: AuthState = {
  user: authRepo.getUser(),
  accessToken: authRepo.getAccessToken(),
  refreshToken: authRepo.getRefreshToken(),
};

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setCredentials: (
      state,
      {
        payload,
      }: PayloadAction<{
        user: User;
        accessToken: string;
        refreshToken: string;
      }>,
    ) => {
      state.user = payload.user;
      authRepo.setUser(payload.user);

      state.accessToken = payload.accessToken;
      authRepo.setAccessToken(payload.accessToken);

      state.refreshToken = payload.refreshToken;
      authRepo.setRefreshToken(payload.refreshToken);
    },

    clearCredentials: (state) => {
      state.user = null;
      authRepo.removeUser();

      state.accessToken = null;
      authRepo.removeAccessToken();

      state.refreshToken = null;
      authRepo.removeRefreshToken();
    },
  },
});

export const authReducer = authSlice.reducer;
export const authActions = authSlice.actions;
