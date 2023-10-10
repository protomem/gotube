import LoginButton from "@/feature/login-button";
import LogoutButton from "@/feature/logout-button";
import { Searcher } from "@/feature/searcher";
import {
  selectIsLoggedIn,
  selectUser,
} from "@/feature/store/auth/auth.selectors";
import { useAppSelector } from "@/feature/store/hooks";
import Title from "@/feature/title";
import { ProfileMenu } from "@/shared/profile-menu";
import { Add } from "@mui/icons-material";
import { Box, IconButton } from "@mui/joy";

export default function AppBar() {
  const isLoggedIn = useAppSelector(selectIsLoggedIn);
  const user = useAppSelector(selectUser);

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

      <Searcher />

      {isLoggedIn && user !== null ? (
        <ProfileMenu
          leftEdge={
            <IconButton>
              <Add />
            </IconButton>
          }
          user={user}
          rightEdge={<LogoutButton />}
        />
      ) : (
        <LoginButton />
      )}
    </Box>
  );
}
