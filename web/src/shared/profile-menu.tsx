import React from "react";
import { User } from "@/entities/models";
import { Box, IconButton } from "@mui/joy";
import Avatar from "@/shared/avatar";

export interface ProfileMenuProps {
  user: User;
  rightEdge?: React.ReactNode;
  leftEdge?: React.ReactNode;
}

export function ProfileMenu({ user, rightEdge, leftEdge }: ProfileMenuProps) {
  return (
    <Box style={{ display: "flex", flexDirection: "row", gap: 20 }}>
      {leftEdge}

      <IconButton>
        <Avatar user={user} />
      </IconButton>

      {rightEdge}
    </Box>
  );
}
