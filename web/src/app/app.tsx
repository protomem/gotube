import QueryProvider from "./query-provider";
import ThemeProvider from "./theme-provider";
import RouteProvider from "./route-provider";
import SideBarStateProvider from "./side-bar-state-provider";
import AuthProvider from "./auth-provider";

const App = () => {
  return (
    <QueryProvider>
      <ThemeProvider>
        <AuthProvider>
          <SideBarStateProvider>
            <RouteProvider />
          </SideBarStateProvider>
        </AuthProvider>
      </ThemeProvider>
    </QueryProvider>
  );
};

export default App;
