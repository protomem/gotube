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

export interface SubscribeRequest {
  toUserNickname: string;
  accessToken: string;
}

export interface UnsubscribeRequest {
  toUserNickname: string;
  accessToken: string;
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

  async subscribe(request: SubscribeRequest) {
    await apiClient.post(`/users/${request.toUserNickname}/subs/`, null, {
      headers: {
        Authorization: "Bearer " + request.accessToken,
      },
    });
  },

  async unsubscribe(request: UnsubscribeRequest) {
    await apiClient.delete(`/users/${request.toUserNickname}/subs/`, {
      headers: {
        Authorization: "Bearer " + request.accessToken,
      },
    });
  },
};
