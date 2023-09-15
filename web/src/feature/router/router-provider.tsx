import {
  createBrowserRouter,
  RouterProvider as ReactRouterProvider,
} from "react-router-dom";
import NotFound from "@/pages/not-found";
import Home from "@/pages/home";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Home />,
  },
  {
    path: "*",
    element: <NotFound />,
  },
]);

export default function RouterProvider() {
  return <ReactRouterProvider router={router} />;
}
