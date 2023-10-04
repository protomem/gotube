import { PayloadAction, createSlice } from "@reduxjs/toolkit";
import { User } from "@/entities/models";
import { authRepo } from "@/entities/auth.repository";

export interface AuthState {
  user: User | null;
  accessToken: string | null;
}

const initialState: AuthState = {
  user: null,
  accessToken: null,
};

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setCredentials: (
      state,
      { payload }: PayloadAction<{ user: User; accessToken: string }>,
    ) => {
      state.user = payload.user;
      authRepo.setUser(payload.user);

      state.accessToken = payload.accessToken;
      authRepo.setAccessToken(payload.accessToken);
    },

    clearCredentials: (state) => {
      state.user = null;
      authRepo.removeUser();

      state.accessToken = null;
      authRepo.removeAccessToken();
    },
  },
});

export const authReducer = authSlice.reducer;
export const authActions = authSlice.actions;
