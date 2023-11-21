import React from "react";

import { useSocketActions, useSocketStore } from "~/state/socket";
import { wsEndpoint, importMovies } from "~/api";

export const Import: React.FC = () => {
  const { messages } = useSocketStore();
  const { connect } = useSocketActions();

  React.useEffect(() => {
    connect(wsEndpoint);
  }, [connect]);

  const onClick = () => {
    console.log("click");
    importMovies();
  };

  return (
    <div className="overflow-auto p-2 text-neutral-300 text-sm bg-slate-800 h-[64rem]">
      {messages.map((line) => (
        <p className="mb0 whitespace-nowrap">{line}</p>
      ))}
      <button
        onClick={onClick}
        className="absolute top10 right-7 bg-blue-700 text-white px-4 py-1 mb-2"
      >
        RUN
      </button>
    </div>
  );
};
