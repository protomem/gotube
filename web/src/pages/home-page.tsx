import { useSearchParams } from "next/navigation";

import { MainLayout } from "@/widgets/layouts/main-layout";
import { AppBar } from "@/widgets/app-bar";
import { SideBar, Navigates } from "@/widgets/side-bar";

export function HomePage() {
  const searchParams = useSearchParams();

  let selectedNav = Navigates.New;
  if (searchParams && searchParams.has("nav")) {
    switch (searchParams.get("nav")) {
      case Navigates.New:
        selectedNav = Navigates.New;
        break;
      case Navigates.Popular:
        selectedNav = Navigates.Popular;
        break;
      default:
        break;
    }
  }

  return (
    <MainLayout appbar=<AppBar /> sidebar=<SideBar selectedNav={selectedNav} />>
      <h1 className="text-4xl">GoTube</h1>
    </MainLayout>
  );
}
