import QueryProvider from "./query-provider";
import ThemeProvider from "./theme-provider";
import RouteProvider from "./route-provider";

export const App = () => {
  return (
    <QueryProvider>
      <ThemeProvider>
        <RouteProvider />
      </ThemeProvider>
    </QueryProvider>
  );
};

export default App;
