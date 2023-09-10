"use client";

import { useSearchParams } from "next/navigation";

export default function Page() {
  const searchParams = useSearchParams();

  const query = searchParams.get("q");

  return <div>Search: {query}</div>;
}
