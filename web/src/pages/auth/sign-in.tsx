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
} from "@chakra-ui/react";

export default function SignIn() {
  return (
    <SingleObjectLayout>
      <Card w="sm">
        <CardHeader>
          <Heading textAlign="center" size="md">
            Sign In
          </Heading>
        </CardHeader>

        <CardBody>
          <Formik
            initialValues={{ email: "", password: "" }}
            onSubmit={() => {}}
          >
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
    </SingleObjectLayout>
  );
}
