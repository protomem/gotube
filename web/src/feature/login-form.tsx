import React from "react";
import { Box, Button, FormControl, FormLabel, Input } from "@mui/joy";
import { useAppDispatch } from "@/feature/store/hooks";
import { useNavigate } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import { authService } from "@/entities/auth.service";
import { authActions } from "@/feature/store/auth/auth.slice";

interface FormElements extends HTMLFormControlsCollection {
  email: HTMLInputElement;
  password: HTMLInputElement;
}

interface LoginFormElements extends HTMLFormElement {
  readonly elements: FormElements;
}

export default function LoginForm() {
  const dispatch = useAppDispatch();
  const nav = useNavigate();

  const mutation = useMutation({
    mutationFn: authService.login,
    onSuccess: (data) => {
      dispatch(authActions.setCredentials(data));

      nav("/", { replace: true });
    },
  });

  const handleSubmit = (e: React.FormEvent<LoginFormElements>) => {
    e.preventDefault();

    mutation.mutate({
      email: e.currentTarget.elements.email.value,
      password: e.currentTarget.elements.password.value,
    });

    e.currentTarget.reset();
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
