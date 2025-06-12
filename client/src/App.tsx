import { createBrowserRouter, RouterProvider } from "react-router-dom";
import RootLayout from "./pages/RootLayout";
import Home from "./pages/Home";
import AboutUs from "./pages/ContactUs";
import AdminSignup from "./pages/AdminSignup";
import AdminLogin from "./components/AdminLogin";

import ProtectedRoutes from "./ProtectedRoutes";
import SecretProtectedPage from "./pages/SecretProtectedPage";
import { Navigate } from "react-router-dom";

const router = createBrowserRouter([
  {
    path: "/",
    element: <RootLayout />,
    children: [
      // {
      //   path: "/", //previously it worked as there was no sub child routes but here there is which causes confusion   //didn't navigate to the home by default //Case 2: With Nested Routes — path: "/" becomes ambiguous aka confusing  //Yes — this renders <Home /> when the user visits /, but it does not redirect the user to /home. //Yes — this renders <Home /> when the user visits /, but it does not redirect the user to /home. //The user sees the Home page, but the URL doesn't change to /home
      //   element: <Home />,
      // },
      {
        index: true, // This means: match when parent path "/" is hit with no subpath //index: true tells React Router: "When the path is exactly /, load this element." //This means: “When user visits the index (empty path "") of this parent (which is /), do this.” ->✅ Equivalent to saying:If path is exactly /, then navigate to /home.
        element: <Navigate to="/home" replace />, //to make it only load after login see the last ques in -> notion page: Final protected route: confusing ambigious when routing react router confusion
      },

      {
        path: "/about",
        element: <AboutUs />,
      },
      {
        path: "/home",
        element: <Home />,
      },
    ],
  },
  {
    path: "/AdminSignUp",
    element: <AdminSignup />,
  },
  {
    path: "/AdminLogin",
    element: <AdminLogin />,
  },
  {
    path: "/secretPage",
    element: (
      <ProtectedRoutes>
        <SecretProtectedPage />
      </ProtectedRoutes>
    ) /* yo bhitra wrap gareko child automatically like prop pass bhayedincha as children = <SecretProtectedPage> jasari */,
  },
]);

const App = () => {
  return <RouterProvider router={router} />;
};

export default App;
