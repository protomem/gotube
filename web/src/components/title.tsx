import React from "react";
import { ROUTES } from "@/lib/routes";

import NextLink from "next/link";
import { Heading } from "@chakra-ui/react";

const Title = () => {
  return (
    <NextLink href={ROUTES.HOME}>
      <Heading size="lg">GoTube</Heading>
    </NextLink>
  );
};

export default Title;
