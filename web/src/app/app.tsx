import QueryProvider from "../providers/query-provider";
import ThemeProvider from "../providers/theme-provider";
import RouteProvider from "../providers/route-provider";
import SideBarStateProvider from "../providers/side-bar-state-provider";
import AuthProvider from "../providers/auth-provider";

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
