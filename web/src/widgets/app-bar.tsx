import LoginButton from "@/feature/login-button";
import LogoutButton from "@/feature/logout-button";
import { selectIsLoggedIn } from "@/feature/store/auth/auth.selectors";
import { useAppSelector } from "@/feature/store/hooks";
import Title from "@/feature/title";
import { Box } from "@mui/joy";

export default function AppBar() {
  const isLoggedIn = useAppSelector(selectIsLoggedIn);

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

      {isLoggedIn ? <LogoutButton /> : <LoginButton />}
    </Box>
  );
}
