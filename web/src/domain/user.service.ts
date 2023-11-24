import { users } from "./fixtures/users";

export const userService = {
  getUser(nickname: string) {
    return users[0];
  },
};
