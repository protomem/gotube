import { UserEntity } from "@/entities/domain/models";

export const users: UserEntity[] = [
  {
    id: "1",
    createdAt: new Date(),
    updatedAt: new Date(),
    nickname: "johndoe",
    email: "johndoe@mail.com",
    isVerified: false,
    avatarPath: "https://picsum.photos/200/300",
    description: "test",
  },
];
