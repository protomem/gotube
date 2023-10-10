import { Avatar as MUIAvatar } from "@mui/joy";
import { User } from "@/entities/models";

export interface AvatarProps {
  user: User;
}

// TODO: Add other colors and image
export default function Avatar({ user }: AvatarProps) {
  return <MUIAvatar size="md">{user.nickname.slice(0, 2)}</MUIAvatar>;
}
