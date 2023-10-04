import React from "react";
import { Box, Button, FormControl, FormLabel, Input } from "@mui/joy";
import { useMutation } from "@tanstack/react-query";
import { authService } from "@/entities/auth.service";
import { useAppDispatch } from "@/feature/store/hooks";
import { authActions } from "./store/auth/auth.slice";
import { useNavigate } from "react-router-dom";

interface FormElements extends HTMLFormControlsCollection {
  nickname: HTMLInputElement;
  email: HTMLInputElement;
  password: HTMLInputElement;
}

interface RegisterFormElements extends HTMLFormElement {
  readonly elements: FormElements;
}

export default function RegisterForm() {
  const dispatch = useAppDispatch();
  const nav = useNavigate();

  const mutation = useMutation({
    mutationFn: authService.register,
    onSuccess: (data) => {
      dispatch(authActions.setCredentials(data));

      nav("/", { replace: true });
    },
  });

  const handleSubmit = (e: React.FormEvent<RegisterFormElements>) => {
    e.preventDefault();

    mutation.mutate({
      nickname: e.currentTarget.elements.nickname.value,
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
          <FormLabel>Nickname</FormLabel>
          <Input type="text" name="nickname" />
        </FormControl>

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
