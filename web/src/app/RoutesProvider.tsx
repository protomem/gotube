import { FC } from "react";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import Home from "../pages/Home";
import SignIn from "../pages/SignIn";
import SignUp from "../pages/SignUp";
import Video from "../pages/Video";
import Channel from "../pages/Channel";
import Search from "../pages/Search";
import Studio from "../pages/Studio";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Home />,
  },
  {
    path: "/auth/sign-in",
    element: <SignIn />,
  },
  {
    path: "/auth/sign-up",
    element: <SignUp />,
  },
  {
    path: "/video/:videoId",
    element: <Video />,
  },
  {
    path: "/channel/:userNickname",
    element: <Channel />,
  },
  {
    path: "/search",
    element: <Search />,
  },
  {
    path: "/studio",
    element: <Studio />,
  },
]);

const RoutesProvider: FC = () => {
  return <RouterProvider router={router} />;
};

export default RoutesProvider;
