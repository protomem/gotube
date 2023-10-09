import { Box } from "@mui/joy";
import NavMenu from "./nav-menu";

export default function SideBar() {
  return (
    <Box
      style={{
        marginTop: 20,
        marginLeft: 10,
        marginRight: 10,
      }}
    >
      <NavMenu />
    </Box>
  );
}
