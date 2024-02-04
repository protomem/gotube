import { useSearchParams } from "next/navigation";
import { Center, Heading } from "@chakra-ui/react";
import MainLayout from "@/layouts/main-layout";
import AppBar from "@/widgets/app-bar";
import SideBar from "@/widgets/side-bar";

export default function Home() {
  const searchParams = useSearchParams();

  const navMentItem = searchParams.get("videos") || "home";

  return (
    <MainLayout
      appBar={<AppBar />}
      sideBar={<SideBar navMenuItemSelected={navMentItem} />}
    >
      <Center>
        <Heading>Home Page: {navMentItem}</Heading>
      </Center>
    </MainLayout>
  );
}
