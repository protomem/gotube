import Link from "next/link";
import { ROUTES } from "@/shared/constants/routes";

export function Title() {
  return (
    <Link href={ROUTES.HOME} className="text-3xl font-bold">
      GoTube
    </Link>
  );
}
