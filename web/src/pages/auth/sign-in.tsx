import { useMutation } from "@tanstack/react-query";
import authService from "@/domain/auth.service";
import { useAuthStore } from "@/domain/stores/auth";
import { useRouter } from "next/navigation";
import NextLink from "next/link";
import SingleObjectLayout from "@/layouts/single-object-layout";
import { Field, Formik } from "formik";
import {
  Button,
  Text,
  Card,
  CardBody,
  CardFooter,
  CardHeader,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Heading,
  Input,
  Link,
  VStack,
  useToast,
} from "@chakra-ui/react";
import ModalSpinner from "@/components/modal-spinner";

type FormValues = {
  email: string;
  password: string;
};

const initialValues: FormValues = {
  email: "",
  password: "",
};

export default function SignIn() {
  const router = useRouter();
  const toast = useToast();
  const login = useAuthStore((state) => state.login);

  const mutation = useMutation({
    mutationFn: authService.signIn,
    onSuccess: ({ data }) => {
      login(data.user, data.accessToken, data.refreshToken);
      router.push("/");
    },
    onError: (error) => {
      toast({
        title: "Sign In Failed",
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
            Sign In
          </Heading>
        </CardHeader>

        <CardBody>
          <Formik initialValues={initialValues} onSubmit={handleSubmit}>
            {({ handleSubmit, errors, touched }) => (
              <form onSubmit={handleSubmit}>
                <VStack spacing={4} align="center">
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
                    Sign In
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

          <Link as={NextLink} href="/auth/sign-up">
            <Text>Account not created yet?</Text>
          </Link>
        </CardFooter>
      </Card>

      <ModalSpinner isOpen={mutation.isPending} />
    </SingleObjectLayout>
  );
}
