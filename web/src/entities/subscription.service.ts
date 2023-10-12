import { Subscription } from "@/entities/models";
import { apiClient } from "@/entities/api";

export interface GetSubscriptionsRequest {
  userNickname: string;
  accessToken: string;
}

export interface GetSubscriptionsResponse {
  subscriptions: Subscription[];
}

export interface GetStatisticsRequest {
  userNickname: string;
}

export interface GetStatisticsResponse {
  subscriptions: number;
  subscribers: number;
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

  async getStatistics(request: GetStatisticsRequest) {
    const response = await apiClient.get<GetStatisticsResponse>(
      `/users/${request.userNickname}/subs/stats`,
    );
    return response.data;
  },
};
