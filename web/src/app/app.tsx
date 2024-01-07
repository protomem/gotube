import QueryProvider from "./query-provider";
import ThemeProvider from "./theme-provider";
import RouteProvider from "./route-provider";
import SideBarStateProvider from "./side-bar-state-provider";

const App = () => {
  return (
    <QueryProvider>
      <ThemeProvider>
        <SideBarStateProvider>
          <RouteProvider />
        </SideBarStateProvider>
      </ThemeProvider>
    </QueryProvider>
  );
};

export default App;
