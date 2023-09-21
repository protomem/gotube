import GoHomeButton from "@/feature/go-home-button";
import AuthForm from "@/widgets/auth-form";
import { Box } from "@mui/joy";

export default function Auth() {
  return (
    <Box
      style={{
        height: "100vh",
        display: "flex",
        flexDirection: "row",
        alignItems: "start",
        justifyContent: "space-between",
      }}
    >
      <Box
        style={{
          flex: 1,
        }}
      >
        <Box
          style={{
            marginTop: "2em",
            marginLeft: "2em",
          }}
        >
          <GoHomeButton />
        </Box>
      </Box>

      <Box
        style={{
          height: "100%",
          flex: 2,
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <AuthForm />
      </Box>

      <Box
        style={{
          flex: 1,
        }}
      />
    </Box>
  );
}
