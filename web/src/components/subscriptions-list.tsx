import { useAuth } from "../providers/auth-provider";
import { subscriptionService } from "../domain/subscription.service";
import { useQuery } from "@tanstack/react-query";
import { Link as RouterLink } from "react-router-dom";
import { Avatar, Link, List, ListItem, Text } from "@chakra-ui/react";

type Props = {
  type?: "expanded" | "minimal";
};

const SubscriptionList = ({ type }: Props) => {
  type = type || "expanded";

  const { currentUser } = useAuth();
  const { data: subscriptions } = useQuery({
    queryKey: ["subscriptions", currentUser?.id],
    queryFn: async () =>
      await subscriptionService.getUserSubscriptions({
        userNickname: currentUser?.nickname || "",
      }),
    select: (data) => data.data.subscriptions,
  });

  return (
    <List>
      {subscriptions?.map((sub) => (
        <ListItem key={sub.id}>
          <Link
            as={RouterLink}
            to={`/profile/${sub.nickname}`}
            _hover={{ textDecoration: "none", bg: "whiteAlpha.200" }}
            rounded="lg"
            py="2"
            pl="2"
            display="flex"
            flexDirection="row"
            alignItems="center"
            gap="2"
          >
            <Avatar name={sub.nickname} src={sub.avatarPath} size="sm" />
            {type === "expanded" && <Text fontSize="lg">{sub.nickname}</Text>}
          </Link>
        </ListItem>
      ))}
    </List>
  );
};

export default SubscriptionList;
