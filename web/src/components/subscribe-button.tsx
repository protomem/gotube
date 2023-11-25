import React from "react";

import { Button } from "@chakra-ui/react";

type SubscribeButtonProps = {
  isSubscribed: boolean;
  onSubscribe: () => void;
  onUnsubscribe: () => void;
};

const SubscribeButton = ({
  isSubscribed,
  onSubscribe,
  onUnsubscribe,
}: SubscribeButtonProps) => {
  return (
    <Button onClick={isSubscribed ? onUnsubscribe : onSubscribe}>
      {isSubscribed ? "Unsubscribe" : "Subscribe"}
    </Button>
  );
};

export default SubscribeButton;
