// import React from "react";

import { Outlet } from "react-router-dom";

import Header from "./shared/header/header";

function App() {
  return (
    <div className="container mx-auto">
      <header>
        <Header />
      </header>
      <main>
        <Outlet />
      </main>
      <footer>footer</footer>
    </div>
  );
}

export default App;
