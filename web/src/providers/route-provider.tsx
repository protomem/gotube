import { createBrowserRouter, RouterProvider } from "react-router-dom";
import HomePage from "../pages/home-page";
import AuthPage from "../pages/auth-page";

const routes = createBrowserRouter([
  {
    path: "/",
    element: <HomePage />,
  },
  {
    path: "/auth",
    element: <AuthPage />,
  },
]);

const RouteProvider = () => {
  return <RouterProvider router={routes} />;
};

export default RouteProvider;
