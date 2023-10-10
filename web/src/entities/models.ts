interface BaseModel {
  id: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface User extends BaseModel {
  nickname: string;
  email: string;
  isVerified: boolean;
  avatarPath: string;
  description: string;
}

export interface Subscription extends BaseModel {
  fromUser: User;
  toUser: User;
}