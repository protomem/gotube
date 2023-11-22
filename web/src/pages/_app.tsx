import "@/styles/globals.css";
import type { AppProps } from "next/app";
import QueryProvider from "@/components/providers/query-provider";
import ThemeProvider from "@/components/providers/theme-provider";

export default function App({ Component, pageProps }: AppProps) {
  return (
    <QueryProvider>
      <ThemeProvider>
        <Component {...pageProps} />
      </ThemeProvider>
    </QueryProvider>
  );
}
