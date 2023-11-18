import { MainLayout } from "@/widgets/layouts/main-layout";
import { AppBar } from "@/widgets/app-bar";
import { SideBar, Navigates } from "@/widgets/side-bar";

export function HomePage() {
  return (
    <MainLayout
      appbar=<AppBar />
      sidebar=<SideBar selectedNav={Navigates.New} />
    >
      <h1 className="text-4xl">GoTube</h1>
    </MainLayout>
  );
}
