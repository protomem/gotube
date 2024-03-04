import { FC } from "react";
import QueryProvider from "./QueryProvider";
import ThemeProvider from "./ThemeProvider";
import RoutesProvider from "./RoutesProvider";

const App: FC = () => {
  return (
    <QueryProvider>
      <ThemeProvider>
        <RoutesProvider />
      </ThemeProvider>
    </QueryProvider>
  );
};

export default App;
