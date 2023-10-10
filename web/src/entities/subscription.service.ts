import { Subscription } from "@/entities/models";
import { apiClient } from "@/entities/api";

export interface GetSubscriptionsRequest {
  userNickname: string;
  accessToken: string;
}

export interface GetSubscriptionsResponse {
  subscriptions: Subscription[];
}

export const subscriptionService = {
  async getSubscriptions(request: GetSubscriptionsRequest) {
    const response = await apiClient.get<GetSubscriptionsResponse>(
      `/users/${request.userNickname}/subs/`,
      {
        headers: {
          Authorization: `Bearer ${request.accessToken}`,
        },
      },
    );
    return response.data;
  },
};
