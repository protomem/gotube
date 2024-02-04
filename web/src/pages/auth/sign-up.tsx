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
} from "@chakra-ui/react";

export default function SignUp() {
  return (
    <SingleObjectLayout>
      <Card w="sm">
        <CardHeader>
          <Heading textAlign="center" size="md">
            Sign Up
          </Heading>
        </CardHeader>

        <CardBody>
          <Formik
            initialValues={{ nickname: "", email: "", password: "" }}
            onSubmit={() => { }}
          >
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
    </SingleObjectLayout>
  );
}
