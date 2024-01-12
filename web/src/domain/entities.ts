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

export interface Video extends BaseEntity {
  title: string;
  description: string;
  thumbnailPath: string;
  videoPath: string;
  author: User;
  views: number;
  isPublic: boolean;
}
