import { Link as RouterLink } from "react-router-dom";
import { Link } from "@mui/joy";

export default function Title() {
  return (
    <Link
      component={RouterLink}
      to="/"
      underline="none"
      style={{ fontSize: "1.7vw", fontWeight: "bold" }}
    >
      GoTube
    </Link>
  );
}
