import {
  createBrowserRouter,
  RouterProvider as ReactRouterProvider,
} from "react-router-dom";
import NotFound from "@/pages/not-found";
import Home from "@/pages/home";
import Auth from "@/pages/auth";
import { Profile } from "@/pages/profile";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Home />,
  },
  {
    path: "/auth",
    element: <Auth />,
  },
  {
    path: "/profile/:userNickname",
    element: <Profile />,
  },
  {
    path: "*",
    element: <NotFound />,
  },
]);

export default function RouterProvider() {
  return <ReactRouterProvider router={router} />;
}
