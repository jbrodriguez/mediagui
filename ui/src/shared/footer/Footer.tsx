import React from "react";

export const Footer: React.FC = () => {
  return (
    <section className="flex flex-row items-center justify-between bg-neutral-100 p-2 mt-4 mb-4">
      <div>
        <span className="text-red-600 mr-1">Copyright &copy;</span>
        <a href="https://jbrio.net/" className="text-sky-700">
          Juan B. Rodriguez
        </a>
      </div>
      <div className="text-red-600">
        <span>mediaGUI &nbsp;</span>
        <span>v1.0</span>
      </div>
      <div className="flex flex-row items-center">
        <a
          className="flex items-center"
          href="https://www.themoviedb.org/"
          title="themoviedb.org"
          target="_blank"
        >
          <img src="/img/tmdb.png" alt="Logo for tmdb" className="w-10 mr-4" />
        </a>

        <a
          className="flex items-center"
          href="https://jbrio.net/"
          title="jbrio.net"
          target="_blank"
        >
          <img src="/img/logo.png" alt="Logo for jbrio.net" className="w-10" />
        </a>
      </div>
    </section>
  );
};
