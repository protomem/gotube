import { useSearchParams } from "next/navigation";
import { AppBar } from "@/widgets/app-bar";
import { MainLayout } from "@/widgets/layouts/main-layout";
import { SideBar } from "@/widgets/side-bar";

export function SearchPage() {
  let query = "";

  const searchParams = useSearchParams();
  if (searchParams !== null && searchParams.has("q")) {
    query = searchParams.get("q") || "";
  }

  return (
    <MainLayout appbar=<AppBar /> sidebar=<SideBar />>
      <h1>Search: {query}</h1>
    </MainLayout>
  );
}
