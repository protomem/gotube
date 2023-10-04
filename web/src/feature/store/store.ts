import { configureStore } from "@reduxjs/toolkit";
import { authReducer } from "@/feature/store/auth/auth.slice";

export const store = configureStore({
  reducer: {
    auth: authReducer,
  },
  devTools: true,
  middleware: (getDefaultMiddleware) => getDefaultMiddleware(),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
