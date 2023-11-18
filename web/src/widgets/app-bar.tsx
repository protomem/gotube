import { ModeToggle } from "@/feature/mode-toggle";
import { ProfileMenu } from "@/feature/profile-menu";
import { Searcher } from "@/feature/searcher";
import { Title } from "@/shared/components/title";

export function AppBar() {
  return (
    <div className="flex flex-row justify-between items-center w-full h-auto px-8">
      <Title />

      <div className="w-[25em]">
        <Searcher />
      </div>

      <div className="flex flex-row gap-4 justify-center items-center">
        <ModeToggle />

        <ProfileMenu />
      </div>
    </div>
  );
}
