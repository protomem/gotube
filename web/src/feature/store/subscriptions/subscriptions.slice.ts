import { Subscription } from "@/entities/models";
import { PayloadAction, createSlice } from "@reduxjs/toolkit";

export interface SubscriptionsState {
  subscriptions: Subscription[];
}

const initialState: SubscriptionsState = {
  subscriptions: [],
};

export const subscriptionsSlice = createSlice({
  name: "subscriptions",
  initialState,
  reducers: {
    setSubscriptions: (state, { payload }: PayloadAction<Subscription[]>) => {
      state.subscriptions = payload;
    },

    addSubscription: (state, { payload }: PayloadAction<Subscription>) => {
      state.subscriptions.push(payload);
    },

    removeSubscription: (state, { payload }: PayloadAction<string>) => {
      state.subscriptions = state.subscriptions.filter((s) => s.id !== payload);
    },
  },
});

export const subscriptionsReducer = subscriptionsSlice.reducer;
export const subscriptionsActions = subscriptionsSlice.actions;
