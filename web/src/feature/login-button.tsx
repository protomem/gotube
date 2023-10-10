import { Button, Typography } from "@mui/joy";
import { useNavigate } from "react-router-dom";

export default function LoginButton() {
  const nav = useNavigate();

  const handleClick = (e: React.MouseEvent) => {
    e.preventDefault();
    nav("/auth", { replace: true });
  };

  return (
    <Button onClick={handleClick}>
      <Typography>login</Typography>
    </Button>
  );
}