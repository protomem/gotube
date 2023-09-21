import { useState } from "react";
import LoginForm from "@/feature/login-form";
import RegisterForm from "@/feature/register-form";
import { Box, Button, ButtonGroup } from "@mui/joy";

enum Forms {
  Login,
  Register,
}

export default function AuthForm() {
  const [selectedForm, setSelectedForm] = useState(Forms.Login);

  return (
    <Box
      style={{
        maxWidth: "25em",
        maxHeight: "30em",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        gap: "2em",
      }}
    >
      <ButtonGroup
        style={{
          width: "100%",
          display: "flex",
          flexDirection: "row",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <Button
          variant={selectedForm === Forms.Login ? "solid" : "outlined"}
          onClick={() => setSelectedForm(Forms.Login)}
        >
          Login
        </Button>

        <Button
          variant={selectedForm === Forms.Register ? "solid" : "outlined"}
          onClick={() => setSelectedForm(Forms.Register)}
        >
          Register
        </Button>
      </ButtonGroup>

      {selectedForm === Forms.Register ? <RegisterForm /> : <LoginForm />}
    </Box>
  );
}
