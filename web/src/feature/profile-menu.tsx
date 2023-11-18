import { Button } from "@/shared/ui/button";
import Link from "next/link";

export function ProfileMenu() {
  return (
    <Button asChild className="text-md">
      <Link href="/not-found">login</Link>
    </Button>
  );
}
