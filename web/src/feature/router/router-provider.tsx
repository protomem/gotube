import {
  createBrowserRouter,
  RouterProvider as ReactRouterProvider,
} from "react-router-dom";
import NotFound from "@/pages/not-found";

const router = createBrowserRouter([
  {
    path: "*",
    element: <NotFound />,
  },
]);

export default function RouterProvider() {
  return <ReactRouterProvider router={router} />;
}
