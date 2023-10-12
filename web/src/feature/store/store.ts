import { configureStore } from "@reduxjs/toolkit";
import { authReducer } from "@/feature/store/auth/auth.slice";
import { navReducer } from "@/feature/store/nav/nav.slice";
import { subscriptionsReducer } from "@/feature/store/subscriptions/subscriptions.slice";

export const store = configureStore({
  reducer: {
    auth: authReducer,
    nav: navReducer,
    subscriptions: subscriptionsReducer,
  },
  devTools: true,
  middleware: (getDefaultMiddleware) => getDefaultMiddleware(),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
