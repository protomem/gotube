import { Button, ButtonGroup } from "@chakra-ui/react";
import React from "react";

const RatingButtons = () => {
  return (
    <ButtonGroup isAttached colorScheme="teal">
      <Button variant="solid" borderRadius="full">
        {"2323 likes"}
      </Button>
      <Button variant="outline" borderRadius="full">
        {"323 dislikes"}
      </Button>
    </ButtonGroup>
  );
};

export default RatingButtons;
