import Title from "@/feature/title";
import { Box } from "@mui/joy";

export default function AppBar() {
  return (
    <Box
      style={{
        width: "100%",
        marginLeft: 20,
        marginRight: 20,
      }}
    >
      <Title />
    </Box>
  );
}
