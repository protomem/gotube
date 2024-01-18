import { SubmitHandler, useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../providers/auth-provider";
import { authService } from "../domain/auth.service";
import { useMutation } from "@tanstack/react-query";
import {
  Button,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Input,
  Text,
  useToast,
} from "@chakra-ui/react";

type FormState = {
  nickname: string;
  email: string;
  password: string;
};

const RegisterForm = () => {
  const nav = useNavigate();
  const toast = useToast();

  const {
    handleSubmit,
    register,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<FormState>();
  const { login } = useAuth();

  const registerMutation = useMutation({
    mutationFn: authService.register,
    onSuccess: (res) => {
      const data = res.data;
      login(data.user, data.accessToken, data.refreshToken);

      reset();
      nav("/", { replace: true });
    },
    onError: (error) => {
      toast({
        title: "Registration failed",
        description: error.message,
        status: "error",
        duration: 3000,
        isClosable: true,
      });

      reset();
    },
  });

  const onSubmit: SubmitHandler<FormState> = (values) => {
    registerMutation.mutate(values);
  };

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      style={{
        display: "flex",
        flexDirection: "column",
        gap: "1rem",
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      <FormControl isInvalid={!!errors.nickname} isRequired>
        <FormLabel htmlFor="nickname">Nickname</FormLabel>
        <Input
          id="nickname"
          placeholder="Nickname"
          {...register("nickname", {
            required: "This is required",
            minLength: { value: 4, message: "Minimum length should be 4" },
          })}
          type="nickname"
        />
        <FormErrorMessage>
          <Text>{errors.nickname && errors.nickname.message}</Text>
        </FormErrorMessage>
      </FormControl>

      <FormControl isInvalid={!!errors.email} isRequired>
        <FormLabel htmlFor="email">Email</FormLabel>
        <Input
          id="email"
          placeholder="Email"
          {...register("email", {
            required: "This is required",
            minLength: { value: 4, message: "Minimum length should be 4" },
          })}
          type="email"
        />
        <FormErrorMessage>
          <Text>{errors.email && errors.email.message}</Text>
        </FormErrorMessage>
      </FormControl>

      <FormControl isInvalid={!!errors.password} isRequired>
        <FormLabel htmlFor="password">Password</FormLabel>
        <Input
          id="password"
          placeholder="Password"
          {...register("password", {
            required: "This is required",
            minLength: { value: 4, message: "Minimum length should be 4" },
          })}
          type="password"
        />
        <FormErrorMessage>
          <Text>{errors.password && errors.password.message}</Text>
        </FormErrorMessage>
      </FormControl>
      <Button mt={4} colorScheme="teal" isLoading={isSubmitting} type="submit">
        Submit
      </Button>
    </form>
  );
};

export default RegisterForm;