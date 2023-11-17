import { useRouter } from "next/router";

import { ROUTES } from "@/shared/constants/routes";
import { Button } from "@/shared/ui/button";

export function ProfileMenu() {
  const route = useRouter();

  return (
    <Button
      className="text-md"
      onClick={() => {
        route.push(ROUTES.LOGIN);
      }}
    >
      login
    </Button>
  );
}
