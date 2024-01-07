interface BaseEntity {
  id: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface User extends BaseEntity {
  nickname: string;
  email: string;
  avatarPath: string;
  description: string;
}
