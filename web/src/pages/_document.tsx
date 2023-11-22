import { ThemeScript } from "@/components/providers/theme-provider";
import { Html, Head, Main, NextScript } from "next/document";

export default function Document() {
  return (
    <Html lang="en">
      <Head />
      <body>
        <ThemeScript />
        <Main />
        <NextScript />
      </body>
    </Html>
  );
}
