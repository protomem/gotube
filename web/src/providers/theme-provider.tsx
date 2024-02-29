import { ReactNode } from "react";
import { ChakraProvider, extendTheme } from "@chakra-ui/react";

const theme = extendTheme({
  fonts: {
    heading: "var(--font-rubik)",
    body: "var(--font-rubik)",
  },
  config: {
    initialColorMode: "dark",
    useSystemColorMode: false,
  },
});

interface Props {
  children: ReactNode;
}

export default function ThemeProvider({ children }: Props) {
  return <ChakraProvider theme={theme}>{children}</ChakraProvider>;
}
