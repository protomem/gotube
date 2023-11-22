import React from "react";

import {
  ChakraProvider,
  ColorModeScript,
  extendTheme,
  type ThemeConfig,
} from "@chakra-ui/react";

const config: ThemeConfig = {
  initialColorMode: "dark",
  useSystemColorMode: false,
};

export const theme = extendTheme({
  config,
});

type ThemeProviderProps = {
  children: React.ReactNode;
};

const ThemeProvider = ({ children }: ThemeProviderProps) => {
  return <ChakraProvider theme={theme}>{children}</ChakraProvider>;
};

export const ThemeScript = () => {
  return <ColorModeScript initialColorMode={theme.config.initialColorMode} />;
};

export default ThemeProvider;
