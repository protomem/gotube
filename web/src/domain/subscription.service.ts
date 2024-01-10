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
    const accessToken = localStorage.getItem("accessToken");

    return await apiClient.get<GetUserSubscriptionsResponse>(
      `/profile/${userNickname}/subs`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      },
    );
  },
};
