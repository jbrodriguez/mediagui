import React from "react";

import type { Movie } from "~/types";
import { hourMinute } from "~/lib/hour-minute";

interface Props {
  movie: Movie;
}

export const CoverScreen: React.FC<Props> = ({ movie }) => {
  return (
    <div className="bg-neutral-100 text-sm">
      <div className="relative overflow-hidden">
        <img src={`/img/p/${movie.cover}`} />
        {movie.count_watched > 0 ? (
          <div className="absolute -left-8 top-5 bg-red-900 -rotate-45 shadow-lg shadow-slate-700">
            <span className="[text-shadow:0px_0px_1px_var(--tw-shadow-color)] shadow-white text-white text-sm px-8 py-1 border">
              watched
            </span>
          </div>
        ) : null}
      </div>
      <div className="p-2">
        <p className="truncate text-sky-700">{movie.title}</p>
        <div className="flex justify-between text-gray-500">
          <span>{movie.year}</span>
          <span>{movie.imdb_rating}</span>
          <span>{hourMinute(movie.runtime)}</span>
        </div>
      </div>
    </div>
  );
};
