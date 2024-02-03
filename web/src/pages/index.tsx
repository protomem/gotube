import { Center, Heading } from "@chakra-ui/react";
import MainLayout from "@/layouts/main-layout";
import AppBar from "@/widgets/app-bar";

export default function Home() {
  return (
    <MainLayout appbar={<AppBar />}>
      <Center>
        <Heading>Home Page</Heading>
      </Center>
    </MainLayout>
  );
}
