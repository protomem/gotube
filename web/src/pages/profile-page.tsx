import { AppBar } from "@/widgets/app-bar";
import { MainLayout } from "@/widgets/layouts/main-layout";
import Error from "next/error";
import { useRouter } from "next/router";

export function ProfilePage() {
  const router = useRouter();
  const { nickname: userNickname } = router.query;
  if (userNickname === undefined) {
    return <Error statusCode={404} title="User not found" />;
  }

  return (
    <MainLayout appbar=<AppBar />>
      <h1>Profile: {userNickname}</h1>
    </MainLayout>
  );
}
