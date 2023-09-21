import React from "react";
import { Box, Button, FormControl, FormLabel, Input } from "@mui/joy";

interface FormElements extends HTMLFormControlsCollection {
  email: HTMLInputElement;
  password: HTMLInputElement;
}

interface LoginFormElements extends HTMLFormElement {
  readonly elements: FormElements;
}

export default function LoginForm() {
  const handleSubmit = (e: React.FormEvent<LoginFormElements>) => {
    e.preventDefault();
  };

  return (
    <Box>
      <form
        onSubmit={handleSubmit}
        style={{
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <FormControl required>
          <FormLabel>Email</FormLabel>
          <Input type="email" name="email" />
        </FormControl>

        <FormControl required>
          <FormLabel>Password</FormLabel>
          <Input type="password" name="password" />
        </FormControl>

        <Button
          type="submit"
          style={{
            marginTop: "1em",
          }}
        >
          Submit
        </Button>
      </form>
    </Box>
  );
}
