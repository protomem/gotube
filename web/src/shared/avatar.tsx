import { Avatar as MUIAvatar } from "@mui/joy";
import { User } from "@/entities/models";

export interface AvatarProps {
  user: User;
  size?: "sm" | "md" | "lg";
}

// TODO: Add other colors and image
export default function Avatar({ user, size }: AvatarProps) {
  if (size === undefined) size = "md";
  return <MUIAvatar size={size}>{user.nickname.slice(0, 2)}</MUIAvatar>;
}
