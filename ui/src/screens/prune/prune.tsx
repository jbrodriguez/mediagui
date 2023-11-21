import React from "react";

import { useSocketActions, useSocketStore } from "~/state/socket";
import { wsEndpoint, pruneMovies } from "~/api";

export const Prune: React.FC = () => {
  const { messages } = useSocketStore();
  const { connect } = useSocketActions();

  React.useEffect(() => {
    connect(wsEndpoint);
  }, [connect]);

  const onClick = () => pruneMovies();

  return (
    <div className="overflow-auto p-2 text-neutral-300 text-sm bg-slate-800 h-[64rem]">
      <button
        onClick={onClick}
        className="absolute top10 right-7 bg-blue-700 text-white px-4 py-1 mb-2"
      >
        RUN
      </button>

      {messages.map((line) => (
        <p className="mb0 whitespace-nowrap">{line}</p>
      ))}
    </div>
  );
};
