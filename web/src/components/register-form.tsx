import { SubmitHandler, useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import {
  Button,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Input,
  Text,
} from "@chakra-ui/react";

type FormState = {
  nickname: string;
  email: string;
  password: string;
};

const RegisterForm = () => {
  const nav = useNavigate();
  const {
    handleSubmit,
    register,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<FormState>();

  const onSubmit: SubmitHandler<FormState> = (values) => {
    console.log(values);
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
