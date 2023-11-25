import React from "react";
import { User } from "@/domain/entities";

import {
  Avatar,
  Box,
  Card,
  CardBody,
  CardHeader,
  Heading,
  SimpleGrid,
  Text,
  VStack,
} from "@chakra-ui/react";

type ProfileCardProps = {
  user: User;
  buttonSubscribe?: React.ReactNode;
};

const ProfileCard = ({ user, buttonSubscribe }: ProfileCardProps) => {
  return (
    <Card align="center" pr={16}>
      <CardHeader>
        <SimpleGrid columns={3}>
          <Box display="flex" flexDirection="column" justifyContent="center">
            {buttonSubscribe}
          </Box>

          <Box display="flex" flexDirection="column" alignItems="center">
            <Avatar name="Dan Abrahmov" />
            <Heading fontSize="xl">{user.nickname}</Heading>
          </Box>

          <Box>
            <VStack>
              <Text>{"12 subscribers"}</Text>
              <Text>{"4 subscriptions"}</Text>
            </VStack>
          </Box>
        </SimpleGrid>
      </CardHeader>

      <CardBody pt={2}>
        <Text>{user.description}</Text>
      </CardBody>
    </Card>
  );
};

export default ProfileCard;
