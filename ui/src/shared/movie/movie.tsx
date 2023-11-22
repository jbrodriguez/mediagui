import React from "react";

import format from "date-fns/format";

import type { Movie } from "~/types";
import { hourMinute } from "~/lib/hour-minute";
import { Icon } from "~/res/icons/icon";
import { Rating } from "~/shared/components/rating";
import { Calendar } from "~/shared/components/calendar";

interface MovieProps {
  index: number;
  item: Movie;
  onFixMovie: ({ index, tmdb_id }: { index: number; tmdb_id: number }) => void;
  onCopyMovie: ({ index, tmdb_id }: { index: number; tmdb_id: number }) => void;
  onRateMovie: ({ index, score }: { index: number; score: number }) => void;
  onWatchedMovie: ({
    index,
    watched,
  }: {
    index: number;
    watched: string;
  }) => void;
}

export const MovieScreen: React.FC<MovieProps> = ({
  index,
  item,
  onFixMovie,
  onCopyMovie,
  onRateMovie,
  onWatchedMovie,
}) => {
  const [value, setValue] = React.useState<number>();
  // const bgImage = `http://localhost:7623/img/b${item.backdrop}`;

  React.useEffect(() => {
    setValue(item.tmdb_id);
  }, [item.tmdb_id]);

  const shows =
    item.all_watched !== ""
      ? item.all_watched
          .split("|")
          .map((show) => format(new Date(show), "MMM dd, yyyy"))
      : [];

  const onValueChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValue(+e.target.value);
  };

  const onFix = () => onFixMovie({ index, tmdb_id: value ?? 0 });
  const onCopy = () => onCopyMovie({ index, tmdb_id: value ?? 0 });
  const onRating = (score: number) => onRateMovie({ index, score });
  const onWatched = (watched: string) => onWatchedMovie({ index, watched });

  return (
    <article
      className="bg-cover bg-center bg-no-repeat pb-2 mb-4"
      style={{
        backgroundImage: `url(/img/b${item.backdrop})`,
      }}
    >
      {/* title, year */}
      <div className="flex justify-between bg-gradient-to-b from-black to-black/65 py-1 px-2">
        <span className="[text-shadow:1px_1px_2px_var(--tw-shadow-color)] shadow-black text-slate-50 font-bold text-2xl">
          {item.title} ({item.year})
        </span>
        <span className="[text-shadow:1px_1px_2px_var(--tw-shadow-color)] shadow-black text-slate-50 font-bold text-2xl">
          {`${hourMinute(item.runtime)}`} | {item.imdb_rating}
        </span>
      </div>

      {/* cover */}
      <div className="pl-2 mt-2">
        <div className="relative overflow-hidden">
          <img
            src={`/img/p${item.cover}`}
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

      <div className="px-2 mt-4">
        <div className="bg-black/25 mt-2 pb-2">
          {/* direct, production_countries */}
          <div className="px-2 flex justify-between">
            <span className="[text-shadow:_0_1px_0_var(--tw-shadow-color),_0_0_1px_var(--tw-shadow-color),_0_1px_5px_var(--tw-shadow-color)] shadow-black/75 text-yellow-400 font-bold">
              {item.director}
            </span>
            <span className="[text-shadow:_0_1px_0_var(--tw-shadow-color),_0_0_1px_var(--tw-shadow-color),_0_1px_5px_var(--tw-shadow-color)] shadow-black/75 text-white font-bold">
              {item.production_countries}
            </span>
          </div>

          {/* actors, genres */}
          <div className="px-2 flex justify-between">
            <span className="[text-shadow:_0_1px_0_var(--tw-shadow-color),_0_0_1px_var(--tw-shadow-color),_0_1px_5px_var(--tw-shadow-color)] shadow-black/75 text-white font-bold">
              {item.actors}
            </span>
            <span className="[text-shadow:_0_1px_0_var(--tw-shadow-color),_0_0_1px_var(--tw-shadow-color),_0_1px_5px_var(--tw-shadow-color)] shadow-black/75 text-white font-bold">
              {item.genres}
            </span>
          </div>

          {/* location */}
          <div className="px-2 mt-4 text-sm">
            <div className="col-span-10">
              <span className="bg-white py-1 px-2 text-sm">
                {item.location}
              </span>
            </div>
          </div>

          {/* id, resolution, added, watched */}
          <div className="flex items-center px-2 mt-4 text-sm">
            <span className=" bg-blue-800 text-slate-50 py-1 px-2">
              {item.id}
            </span>
            <span className=" bg-blue-700 text-slate-50 ml-2 py-1 px-2">
              {item.resolution}
            </span>
            <div className="bg-green-700 text-white ml-2 py-1 px-2 flex items-center">
              <Icon name="plus" size={12} fill="fill-white" />
              <span className="ml-2">
                {format(new Date(item.added), "MMM dd, yyyy")}
              </span>
            </div>
            {item.count_watched > 0 ? (
              <div className="bg-blue-600 text-white ml-2 py-1 px-2 flex items-center">
                <Icon name="binoculars" size={12} fill="fill-white" />
                <span className="ml-2">
                  {format(new Date(item.last_watched), "MMM dd, yyyy")}
                </span>
              </div>
            ) : null}
          </div>

          {/* overview */}
          <div className="px-2 mt-4">
            <span className="[text-shadow:_0_1px_0_var(--tw-shadow-color),_0_0_1px_var(--tw-shadow-color),_0_1px_5px_var(--tw-shadow-color)] shadow-black/75 text-white font-bold">
              {item.overview}
            </span>
          </div>

          {/* tmdb, fix, copy, dup, count_watched, history, score, watched input */}
          <div className="px-2 mt-4 flex justify-between text-sm">
            <div className="flex flex-row">
              <input
                type="number"
                defaultValue={item.tmdb_id}
                value={value}
                onChange={onValueChange}
                className="bg-white text-slate-600 px-2 py-1"
              />
              <button
                className="bg-blue-700 text-white px-4 py-1 ml-2"
                onClick={onFix}
              >
                fix
              </button>
              <button
                className="bg-blue-700 text-white px-2 py-1 ml-2"
                onClick={onCopy}
              >
                copy
              </button>
              <span className="ml-2 flex items-center">
                <span className="text-white">!dup? </span>
                <input type="checkbox" id="dup" className="ml-1" />
              </span>
            </div>

            {item.count_watched > 0 ? (
              <div className="flex items-center">
                <span className="font-bold me-2 px-2.5 py-1 rounded-full bg-green-900 text-green-300">
                  {item.count_watched}
                </span>
                <select
                  className="mr-2 text-slate-600 py-1 px-2 outline-0"
                  value={`${shows[shows.length - 1]}`}
                >
                  {shows.map((show, index) => (
                    <option key={index}>{show}</option>
                  ))}
                </select>
              </div>
            ) : null}

            <div className="flex items-center">
              {item.score > 0 ? (
                <span className="font-bold me-2 px-2.5 py-0.5 rounded bg-red-900 text-red-300">
                  {item.score}
                </span>
              ) : null}
              <Rating rating={item.score} setRating={onRating} />
              <Calendar onChange={onWatched} />
            </div>
          </div>
        </div>
      </div>
    </article>
  );
};
