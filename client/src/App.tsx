import { createBrowserRouter, RouterProvider } from "react-router-dom";
import RootLayout from "./pages/RootLayout";
import UserLogin from "./pages/UserLogin";

const router = createBrowserRouter([
  {
    path: "/",
    element: <RootLayout />,
    children: [
      {
        index: true,
        element: <UserLogin />,
      },
    ],
  },
]);

const App = () => {
  return <RouterProvider router={router}/>;
};

export default App;
