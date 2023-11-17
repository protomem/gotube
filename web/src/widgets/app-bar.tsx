import { ProfileMenu } from "@/feature/profile-menu";
import { Title } from "@/shared/components/title";

export function AppBar() {
  return (
    <div className="flex flex-row justify-between items-center w-full h-auto px-8">
      <Title />

      <ProfileMenu />
    </div>
  );
}
