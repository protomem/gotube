import { FC } from "react";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import Home from "../pages/Home";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Home />,
  },
]);

const RoutesProvider: FC = () => {
  return <RouterProvider router={router} />;
};

export default RoutesProvider;
