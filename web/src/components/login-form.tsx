import { SubmitHandler, useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../app/auth-provider";
import {
  Button,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Input,
  Text,
} from "@chakra-ui/react";

type FormState = {
  email: string;
  password: string;
};

const LoginForm = () => {
  const nav = useNavigate();
  const {
    handleSubmit,
    register,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<FormState>();
  const { login } = useAuth();

  const onSubmit: SubmitHandler<FormState> = (values) => {
    login(
      {
        id: "1",
        createdAt: new Date(),
        updatedAt: new Date(),
        nickname: values.email.split("@")[0],
        email: values.email,
        avatarPath: "",
        description: "",
      },
      "access_some_token",
      "refresh_some_token"
    );

    reset();
    nav("/", { replace: true });
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

export default LoginForm;
