import React from "react";

import format from "date-fns/format";

import type { Movie } from "~/types";
import { hourMinute } from "~/lib/hour-minute";
import { Icon } from "~/res/icons/icon";

// TODO(jbrodriguez): this is a hack, it works for my screen size, but it's not scalable
const maxwidth = 175;

interface MovieProps {
  item: Movie;
}

const Movie: React.FC<MovieProps> = ({ item }) => {
  // const bgImage = `http://localhost:7623/img/b${item.backdrop}`;
  const styles = item.location.length > maxwidth ? "whitespace-nowrap" : "";
  const location =
    item.location.length > maxwidth
      ? "..." + item.location.slice(-maxwidth)
      : item.location;

  return (
    <article
      className="bg-cover bg-center bg-no-repeat pb-2 mb-4"
      style={{
        backgroundImage: `url(http://localhost:7623/img/b${item.backdrop})`,
      }}
    >
      <div className="flex justify-between bg-gradient-to-b from-black to-black/65 py-1 px-2">
        <span className="[text-shadow:1px_1px_2px_var(--tw-shadow-color)] shadow-black text-slate-50 font-bold text-2xl">
          {item.title} ({item.year})
        </span>
        <span className="[text-shadow:1px_1px_2px_var(--tw-shadow-color)] shadow-black text-slate-50 font-bold text-2xl">
          {`${hourMinute(item.runtime)}`} | {item.imdb_rating}
        </span>
      </div>

      <div className="pl-2">
        <div className="relative overflow-hidden">
          <img
            src={`http://localhost:7623/img/p${item.cover}`}
            className="h-96 opacity-75"
            alt={item.title}
          />
          {item.count_watched > 0 ? (
            <div className="absolute -left-8 top-5 bg-red-900 -rotate-45 shadow-lg shadow-slate-700">
              <span className="[text-shadow:0px_0px_1px_var(--tw-shadow-color)] shadow-white text-white text-sm px-8 py-1 border">
                watched
              </span>
            </div>
          ) : null}
        </div>
      </div>

      <div className="px-2">
        <div className="bg-black/25 mt-2 pb-2">
          <div className="px-2 flex justify-between">
            <span className="[text-shadow:_0_1px_0_var(--tw-shadow-color),_0_0_1px_var(--tw-shadow-color),_0_1px_5px_var(--tw-shadow-color)] shadow-black/75 text-yellow-400 font-bold">
              {item.director}
            </span>
            <span className="[text-shadow:_0_1px_0_var(--tw-shadow-color),_0_0_1px_var(--tw-shadow-color),_0_1px_5px_var(--tw-shadow-color)] shadow-black/75 text-white font-bold">
              {item.production_countries}
            </span>
          </div>

          <div className="px-2 flex justify-between">
            <span className="[text-shadow:_0_1px_0_var(--tw-shadow-color),_0_0_1px_var(--tw-shadow-color),_0_1px_5px_var(--tw-shadow-color)] shadow-black/75 text-white font-bold">
              {item.actors}
            </span>
            <span className="[text-shadow:_0_1px_0_var(--tw-shadow-color),_0_0_1px_var(--tw-shadow-color),_0_1px_5px_var(--tw-shadow-color)] shadow-black/75 text-white font-bold">
              {item.genres}
            </span>
          </div>

          <div className="px-2 mt-4">
            <span className="[text-shadow:_0_1px_0_var(--tw-shadow-color),_0_0_1px_var(--tw-shadow-color),_0_1px_5px_var(--tw-shadow-color)] shadow-black/75 text-white font-bold">
              {item.overview}
            </span>
          </div>

          <div className="grid grid-cols-12 gap-0 px-2 mt-2">
            <div className="col-span-11">
              <span className=" bg-blue-800 text-slate-50 text-sm py-1 px-2">
                {item.id}
              </span>
              <span className={`bg-white py-1 px-2 text-sm ${styles}`}>
                {location}
              </span>
              <span className=" bg-blue-800 text-slate-50 text-sm py-1 px-2">
                {item.resolution}
              </span>
            </div>

            <div className="col-span-1 flex items-center justify-end">
              <div className="bg-green-700 text-white text-xs font-bold py-1 px-2 flex items-center ml-2">
                <Icon name="plus" size={12} />
                <span className="ml-2">
                  {format(new Date(item.added), "MMM dd, yyyy")}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </article>
  );
};

export default Movie;
