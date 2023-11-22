import React from "react";

import { useSocketActions, useSocketStore } from "~/state/socket";
import { wsEndpoint, importMovies } from "~/api";

export const Import: React.FC = () => {
  const { messages } = useSocketStore();
  const { connect } = useSocketActions();

  React.useEffect(() => {
    connect(wsEndpoint);
    // return () => {

    // }
  }, [connect]);

  const onClick = () => importMovies();

  return (
    <div className="relative overflow-auto p-2 text-neutral-300 text-sm bg-slate-800 h-[64rem]">
      {messages.map((line) => (
        <p className="mb0 whitespace-nowrap">{line}</p>
      ))}
      <button
        onClick={onClick}
        className="absolute right-2 bg-blue-700 text-white px-4 py-1 mb-2"
      >
        RUN
      </button>
    </div>
  );
};
