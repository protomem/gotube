import { Box, Typography } from "@mui/joy";

export default function NotFound() {
  return (
    <Box
      style={{
        height: "100vh",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
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
    </Box>
  );
}
