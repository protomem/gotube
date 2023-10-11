import { Button, Typography } from "@mui/joy";
import { useAppDispatch, useAppSelector } from "@/feature/store/hooks";
import { useNavigate } from "react-router-dom";
import { authActions } from "@/feature/store/auth/auth.slice";
import { useMutation } from "@tanstack/react-query";
import { authService } from "@/entities/auth.service";
import {
  selectAccessToken,
  selectRefreshToken,
} from "@/feature/store/auth/auth.selectors";

export default function LogoutButton() {
  const dispatch = useAppDispatch();
  const nav = useNavigate();

  const accessToken = useAppSelector(selectAccessToken);
  const refreshToken = useAppSelector(selectRefreshToken);

  const mutation = useMutation({
    mutationFn: authService.logout,
  });

  const handleClick = (e: React.MouseEvent) => {
    e.preventDefault();

    mutation.mutate({
      accessToken: accessToken || "",
      refreshToken: refreshToken || "",
    });

    dispatch(authActions.clearCredentials());
    nav("/", { replace: true });
  };

  return (
    <Button onClick={handleClick} color="danger" variant="soft">
      <Typography>logout</Typography>
    </Button>
  );
}
