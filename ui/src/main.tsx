import React from "react";
import ReactDOM from "react-dom/client";
import {
  createBrowserRouter,
  RouterProvider,
  Navigate,
} from "react-router-dom";

import App from "./App.tsx";
import "./index.css";
import { Movies } from "~/screens/movies/movies.tsx";
import { Import } from "~/screens/import/import.tsx";
import { Prune } from "~/screens/prune/prune.tsx";
import { Duplicates } from "~/screens/duplicates/duplicates.tsx";

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    children: [
      { index: true, element: <Navigate to="/movies" replace /> },
      {
        path: "/movies",
        element: <Movies />,
      },
      {
        path: "/import",
        element: <Import />,
      },
      {
        path: "/prune",
        element: <Prune />,
      },
      {
        path: "/duplicates",
        element: <Duplicates />,
      },
    ],
  },
]);

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
);
