import { RootState } from "@/feature/store/store";

const selectSubscriptionsModule = (state: RootState) => state.subscriptions;

export const selectSubscriptions = (state: RootState) =>
  selectSubscriptionsModule(state).subscriptions;
