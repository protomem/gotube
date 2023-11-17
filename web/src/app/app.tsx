import type { AppProps } from "next/app";
import { lora } from "@/shared/fonts";

import { AppProvider } from "@/app/app-provider";

export function App({ Component, pageProps }: AppProps) {
  return (
    <AppProvider>
      <main className={`${lora.variable} font-sans`}>
        <Component {...pageProps} />
      </main>
    </AppProvider>
  );
}
