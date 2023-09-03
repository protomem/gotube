"use client";

import { Card, CardHeader, CardTitle } from "@/shared/ui/card";

export default function Home() {
  console.log(process.env.NEXT_PUBLIC_APP_URL);

  return (
    <main className="flex flex-coll justify-center align-center mt-8">
      <Card className="w-[30rem] h-[10rem]">
        <CardHeader>
          <CardTitle className="text-3xl text-center m-8">GoTube</CardTitle>
        </CardHeader>
      </Card>
    </main>
  );
}
