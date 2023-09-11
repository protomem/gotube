import React from "react";
import StoreProvider from "@/feature/store/store-provider";
import QueryProvider from "@/feature/query/query-provider";
import ThemeProvider from "@/feature/theme/theme-provider";

interface SettingsProps {
  children: React.ReactNode;
}

export default function Settings({ children }: SettingsProps) {
  return (
    <StoreProvider>
      <QueryProvider>
        <ThemeProvider>{children}</ThemeProvider>
      </QueryProvider>
    </StoreProvider>
  );
}
