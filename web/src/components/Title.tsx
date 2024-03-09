import { FC } from "react";
import { Link as RouterLink } from "react-router-dom";
import { Heading, Link } from "@chakra-ui/react";

const Title: FC = () => {
  return (
    <Link as={RouterLink} to="/" style={{ textDecoration: "none" }}>
      <Heading size="xl">GoTube</Heading>
    </Link>
  );
};

export default Title;
