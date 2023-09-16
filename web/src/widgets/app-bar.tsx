import LoginButton from "@/feature/login-button";
import Title from "@/feature/title";
import { Box } from "@mui/joy";

export default function AppBar() {
  return (
    <Box
      style={{
        width: "100%",
        marginLeft: 50,
        marginRight: 50,
        display: "flex",
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "space-between",
      }}
    >
      <Title />

      <LoginButton />
    </Box>
  );
}
