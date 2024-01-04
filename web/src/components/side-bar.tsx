import { Box } from "@chakra-ui/react";

type Props = {
  type?: "minimal" | "expanded";
};

const SideBar = ({ type }: Props) => {
  type = type || "expanded";

  return (
    <Box width={type === "minimal" ? "6rem" : "16rem"}>
      Side Bar {type === "expanded" && "expanded"}
    </Box>
  );
};

export default SideBar;
