interface BaseEntity {
  id: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface User extends BaseEntity {
  nickname: string;
  email: string;
  isVerified: boolean;
  avatarPath: string;
  description: string;
}

export interface Video extends BaseEntity {
  title: string;
  description: string;
  thumbnailPath: string;
  videoPath: string;
  author: User;
  isPublic: boolean;
  views: number;
}
