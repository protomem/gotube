interface BaseEnity {
  id: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface User extends BaseEnity {
  nickname: string;
  email: string;
  isVerified: boolean;
  avatarPath: string;
  description: string;
}

export interface Video extends BaseEnity {
  title: string;
  description: string;
  thumbnailPath: string;
  videoPath: string;
  author: User;
  isPublic: boolean;
  views: number;
}

export interface Comment extends BaseEnity {
  content: string;
  author: User;
  videoId: string;
}
