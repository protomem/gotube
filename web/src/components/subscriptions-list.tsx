import { useState } from "react";
import { User } from "../domain/entities";
import { Link as RouterLink } from "react-router-dom";
import { Avatar, Link, List, ListItem, Text } from "@chakra-ui/react";

type Props = {
  type?: "expanded" | "minimal";
  user: User;
};

const SubscriptionList = ({ type }: Props) => {
  type = type || "expanded";
  console.log(type);

  const [subscriptions] = useState<User[]>([]);

  // TODO: add tooltip
  // TODO: add query
  return (
    <List>
      {subscriptions.map((sub) => (
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
