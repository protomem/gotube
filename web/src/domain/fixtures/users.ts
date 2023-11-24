import { User } from "../entities";

export const users: User[] = [
  {
    id: "1",
    createdAt: new Date(),
    updatedAt: new Date(),
    nickname: "roman",
    email: "roman@ya.ru",
    isVerified: false,
    avatarPath: "",
    description: "Description for roman",
  },
];
