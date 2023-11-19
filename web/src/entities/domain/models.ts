interface BaseEntity {
  id: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface UserEntity extends BaseEntity {
  nickname: string;
  email: string;
  isVerified: boolean;
  avatarPath: string;
  description: string;
}

export interface VideoEntity extends BaseEntity {
  title: string;
  description: string;
  thumbnailPath: string;
  videoPath: string;
  author: UserEntity;
  isPublic: boolean;
  views: number;
}
