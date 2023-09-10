import axios from "axios";

export const apiClient = axios.create({
  baseURL: `http://${process.env.NEXT_PUBLIC_APP_URL}/api/v1`,
});
