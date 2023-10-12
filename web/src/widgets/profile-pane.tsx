import { useState } from "react";
import { User } from "@/entities/models";
import Avatar from "@/shared/avatar";
import { Box, Button, Typography } from "@mui/joy";
import { useQuery } from "@tanstack/react-query";
import { subscriptionService } from "@/entities/subscription.service";
import { useAppSelector } from "@/feature/store/hooks";
import { selectIsLoggedIn } from "@/feature/store/auth/auth.selectors";

export interface ProfilePaneProps {
  user: User;
}

export default function ProfilePane({ user }: ProfilePaneProps) {
  const isLoggedIn = useAppSelector(selectIsLoggedIn);

  const [subsStat, setSubsStat] = useState({
    countSubscriptions: 0,
    countSubscribers: 0,
  });

  useQuery({
    queryKey: ["subscriptions", "stats", user.nickname],
    queryFn: async () =>
      subscriptionService.getStatistics({
        userNickname: user.nickname,
      }),
    onSuccess: (data) => {
      setSubsStat({
        countSubscriptions: data.subscriptions,
        countSubscribers: data.subscribers,
      });
    },
  });

  return (
    <Box
      style={{
        display: "flex",
        flexDirection: "column",
        gap: 20,
      }}
    >
      <Box
        style={{
          display: "flex",
          flexDirection: "row",
          justifyContent: "center",
          gap: 10,
        }}
      >
        <Box
          style={{
            flex: 4,
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
            alignItems: "end",
            gap: 10,
          }}
        >
          {isLoggedIn && (
            <Button size="sm" style={{ width: "20%" }}>
              subscribe
            </Button>
          )}
        </Box>

        <Box
          style={{
            flex: 1,
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            gap: 10,
          }}
        >
          <Avatar size="lg" user={user} />
          <Typography level="h3" style={{ textAlign: "center" }}>
            {user.nickname}
          </Typography>
        </Box>

        <Box
          style={{
            flex: 4,
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
            alignItems: "start",
            gap: 10,
          }}
        >
          <Typography
            level="body-sm"
            fontWeight="bold"
            style={{ textAlign: "center" }}
          >
            {subsStat.countSubscriptions} subscriptions
          </Typography>
          <Typography
            level="body-sm"
            fontWeight="bold"
            style={{ textAlign: "center" }}
          >
            {subsStat.countSubscribers} subscribers
          </Typography>
        </Box>
      </Box>

      {user.description !== "" && (
        <Typography level="body-lg" style={{ textAlign: "left" }}>
          {user.description}
        </Typography>
      )}
    </Box>
  );
}
