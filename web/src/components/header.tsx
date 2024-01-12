import { Link as RouterLink } from "react-router-dom";
import { Heading, Link } from "@chakra-ui/react";

const Header = () => {
  return (
    <Link as={RouterLink} to="/" _hover={{ textDecoration: "none" }}>
      <Heading size="lg">GoTube</Heading>
    </Link>
  );
};

export default Header;
