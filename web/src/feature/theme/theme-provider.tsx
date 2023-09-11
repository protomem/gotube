import { CssBaseline, CssVarsProvider } from "@mui/joy";
import React from "react";
import { theme } from "@/feature/theme/theme";

interface ThemeProviderProps {
  children: React.ReactNode;
}

export default function ThemeProvider({ children }: ThemeProviderProps) {
  return (
    <CssVarsProvider theme={theme} defaultMode="dark">
      <CssBaseline />

      {children}
    </CssVarsProvider>
  );
}
