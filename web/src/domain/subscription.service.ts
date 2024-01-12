import { apiClient } from "./api.client";
import { User } from "./entities";

type GetUserSubscriptionsRequest = {
  userNickname: string;
};

type GetUserSubscriptionsResponse = {
  subscriptions: User[];
};

export const subscriptionService = {
  async getUserSubscriptions({ userNickname }: GetUserSubscriptionsRequest) {
    return await apiClient.get<GetUserSubscriptionsResponse>(
      `/profile/${userNickname}/subs`,
    );
  },
};
