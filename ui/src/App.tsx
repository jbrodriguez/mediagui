// import React from "react";

import { Outlet } from "react-router-dom";

import { Header } from "./shared/header/header";
import { Footer } from "./shared/footer/Footer";

function App() {
  return (
    <div className="container mx-auto">
      <header>
        <Header />
      </header>
      <main>
        <Outlet />
      </main>
      <footer>
        <Footer />
      </footer>
    </div>
  );
}

export default App;
