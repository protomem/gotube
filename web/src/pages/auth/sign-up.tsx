import { useRouter } from "next/navigation";
import { useMutation } from "@tanstack/react-query";
import authService from "@/domain/auth.service";
import { useAuthStore } from "@/domain/stores/auth";
import NextLink from "next/link";
import SingleObjectLayout from "@/layouts/single-object-layout";
import { Formik, Field } from "formik";
import {
  Card,
  Button,
  Input,
  CardBody,
  CardHeader,
  CardFooter,
  Link,
  Heading,
  Text,
  VStack,
  FormLabel,
  FormControl,
  FormErrorMessage,
  useToast,
} from "@chakra-ui/react";
import ModalSpinner from "@/components/modal-spinner";

type FormValues = {
  nickname: string;
  email: string;
  password: string;
};

const initialValues: FormValues = {
  nickname: "",
  email: "",
  password: "",
};

export default function SignUp() {
  const router = useRouter();
  const toast = useToast();
  const login = useAuthStore((state) => state.login);

  const mutation = useMutation({
    mutationFn: authService.signUp,
    onSuccess: ({ data }) => {
      login(data.user, data.accessToken, data.refreshToken);
      router.push("/");
    },
    onError: (error) => {
      toast({
        title: "Sign Up Failed",
        description: "Oops, something went wrong. Please try again later.",
        status: "error",
        duration: 4000,
        isClosable: true,
      });
      console.error(`sign-up-error: ${error}`);
    },
  });

  const handleSubmit = (values: FormValues) => {
    mutation.mutate(values);
  };

  return (
    <SingleObjectLayout>
      <Card w="sm">
        <CardHeader>
          <Heading textAlign="center" size="md">
            Sign Up
          </Heading>
        </CardHeader>

        <CardBody>
          <Formik initialValues={initialValues} onSubmit={handleSubmit}>
            {({ handleSubmit, errors, touched }) => (
              <form onSubmit={handleSubmit}>
                <VStack spacing={4} align="center">
                  <FormControl>
                    <FormLabel htmlFor="nickname">Nickname</FormLabel>
                    <Field
                      as={Input}
                      id="nickname"
                      name="nickname"
                      type="nickname"
                      variant="filled"
                      autoComplete="off"
                    />
                  </FormControl>

                  <FormControl>
                    <FormLabel htmlFor="email">Email</FormLabel>
                    <Field
                      as={Input}
                      id="email"
                      name="email"
                      type="email"
                      variant="filled"
                      autoComplete="off"
                    />
                  </FormControl>

                  <FormControl
                    isInvalid={!!errors.password && touched.password}
                  >
                    <FormLabel htmlFor="password">Password</FormLabel>
                    <Field
                      as={Input}
                      id="password"
                      name="password"
                      type="password"
                      variant="filled"
                      validate={(value: string) => {
                        let error;

                        if (value.length < 6) {
                          error = "Password must contain at least 6 characters";
                        }

                        return error;
                      }}
                    />
                    <FormErrorMessage>{errors.password}</FormErrorMessage>
                  </FormControl>

                  <Button type="submit" variant="solid" colorScheme="teal">
                    Sign Up
                  </Button>
                </VStack>
              </form>
            )}
          </Formik>
        </CardBody>

        <CardFooter display="flex" justifyContent="space-between">
          <Link as={NextLink} href="/">
            <Text>Home</Text>
          </Link>

          <Link as={NextLink} href="/auth/sign-in">
            <Text>Account already exists?</Text>
          </Link>
        </CardFooter>
      </Card>

      <ModalSpinner isOpen={mutation.isPending} />
    </SingleObjectLayout>
  );
}
