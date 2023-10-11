import { useState } from "react";
import { List, ListItemButton, ListSubheader, Typography } from "@mui/joy";
import { useAppSelector } from "@/feature/store/hooks";
import {
  selectAccessToken,
  selectUser,
} from "@/feature/store/auth/auth.selectors";
import { Subscription } from "@/entities/models";
import { useQuery } from "@tanstack/react-query";
import { subscriptionService } from "@/entities/subscription.service";
import Avatar from "@/shared/avatar";
import { useNavigate } from "react-router-dom";

export default function SubscriptionsList() {
  const nav = useNavigate();

  const user = useAppSelector(selectUser);
  const accessToken = useAppSelector(selectAccessToken);

  const [subs, setSubs] = useState<Subscription[]>([]);

  useQuery({
    queryKey: ["subscriptions", user?.nickname],
    queryFn: async () =>
      subscriptionService.getSubscriptions({
        userNickname: user?.nickname || "",
        accessToken: accessToken || "",
      }),
    onSuccess: (data) => {
      setSubs(data.subscriptions);
    },
    enabled: user !== null && accessToken !== null,
  });

  return (
    <List style={{ gap: 5 }}>
      <ListSubheader>Subscriptions</ListSubheader>

      {subs.map((s) => (
        <ListItemButton
          key={s.id}
          onClick={() => {
            nav(`/profile/${s.toUser.nickname}`, { replace: true });
          }}
        >
          <Avatar user={s.toUser} />
          <Typography style={{ marginLeft: 10 }}>
            {s.toUser.nickname}
          </Typography>
        </ListItemButton>
      ))}
    </List>
  );
}
