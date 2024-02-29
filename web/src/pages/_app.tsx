import "@/styles/globals.css";
import type { AppProps } from "next/app";
import { fonts } from "@/lib/fonts";
import ThemeProvider from "@/providers/theme-provider";
import QueryProvider from "@/providers/query-provider";
import StoresProvider from "@/providers/stores-provider";

export default function App({ Component, pageProps }: AppProps) {
  return (
    <>
      <style jsx global>
        {`
          :root {
            --font-rubik: ${fonts.rubik.style.fontFamily};
          }
        `}
      </style>

      <QueryProvider>
        <StoresProvider>
          <ThemeProvider>
            <Component {...pageProps} />
          </ThemeProvider>
        </StoresProvider>
      </QueryProvider>
    </>
  );
}
