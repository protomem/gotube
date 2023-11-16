import "@/styles/globals.css";
import type { AppProps } from "next/app";

import { ThemeProvider } from "@/components/settings/theme-provider";
import { QueryProvider } from "@/components/settings/query-provider";

export default function App({ Component, pageProps }: AppProps) {
  return (
    <QueryProvider>
      <ThemeProvider
        attribute="class"
        defaultTheme="dark"
        disableTransitionOnChange
      >
        <Component {...pageProps} />
      </ThemeProvider>
    </QueryProvider>
  );
}
