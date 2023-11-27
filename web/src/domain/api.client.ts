import axios from "axios";

export const apiClient = axios.create({
  baseURL: "/api/v1",
  headers: {
    ContentType: "application/json",
  },
});
