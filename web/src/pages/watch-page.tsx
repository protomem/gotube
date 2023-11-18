import { AppBar } from "@/widgets/app-bar";
import { MainLayout } from "@/widgets/layouts/main-layout";
import Error from "next/error";
import { useRouter } from "next/router";

export function WatchPage() {
  const router = useRouter();
  const { id: videoId } = router.query;
  if (videoId === undefined) {
    return <Error statusCode={404} title="Video not found" />;
  }

  return (
    <MainLayout appbar=<AppBar />>
      <h1>Video: {videoId}</h1>
    </MainLayout>
  );
}
