import GoHomeButton from "@/feature/go-home-button";
import { Box, Typography } from "@mui/joy";

export default function NotFound() {
  return (
    <Box
      style={{
        height: "100vh",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        gap: "18px",
      }}
    >
      <Box
        style={{
          display: "flex",
          flexDirection: "row",
          alignItems: "center",
          justifyContent: "center",
          gap: "18px",
        }}
      >
        <Typography fontSize={25} fontWeight={"bold"} textAlign={"center"}>
          404
        </Typography>
        <Typography fontSize={30}>|</Typography>
        <Typography fontSize={25} textAlign={"center"}>
          Not found
        </Typography>
      </Box>

      <GoHomeButton withArrow={false} />
    </Box>
  );
}
