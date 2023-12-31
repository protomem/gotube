import { createBrowserRouter, RouterProvider } from "react-router-dom";
import HomePage from "../pages/home-page";

const routes = createBrowserRouter([
  {
    path: "/",
    element: <HomePage />,
  },
]);

export const RouteProvider = () => {
  return <RouterProvider router={routes} />;
};

export default RouteProvider;