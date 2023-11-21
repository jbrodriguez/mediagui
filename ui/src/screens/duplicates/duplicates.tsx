import useSWR from "swr";

import {
  getDuplicates,
  fixMovie,
  copyMovie,
  rateMovie,
  watchedMovie,
} from "~/api";
import { MovieScreen } from "~/shared/movie/movie";

export const Duplicates = () => {
  const { data, error, mutate } = useSWR("/duplicates", getDuplicates);

  if (error) return <div>Error</div>;
  if (!data) return <div>No data</div>;

  // console.log("data", data);
  const total = data.total ?? 0;

  const onFixMovie = async ({
    index,
    tmdb_id,
  }: {
    index: number;
    tmdb_id: number;
  }) => {
    // const index = data.items.findIndex((item) => item.id === id);
    const id = data.items[index].id;
    data.items[index] = await fixMovie({ id, tmdb_id });
    mutate({ items: [...data.items], total }, { revalidate: false });
  };

  const onCopyMovie = async ({
    index,
    tmdb_id,
  }: {
    index: number;
    tmdb_id: number;
  }) => {
    // const index = data.items.findIndex((item) => item.id === id);
    const id = data.items[index].id;
    data.items[index] = await copyMovie({ id, tmdb_id });
    mutate({ items: [...data.items], total }, { revalidate: false });
  };

  const onRateMovie = async ({
    index,
    score,
  }: {
    index: number;
    score: number;
  }) => {
    // const index = data.items.findIndex((item) => item.id === id);
    const id = data.items[index].id;
    data.items[index] = await rateMovie({ id, score });
    mutate({ items: [...data.items], total }, { revalidate: false });
  };

  const onWatchedMovie = async ({
    index,
    watched,
  }: {
    index: number;
    watched: string;
  }) => {
    // const index = data.items.findIndex((item) => item.id === id);
    const id = data.items[index].id;
    data.items[index] = await watchedMovie({ id, watched });
    mutate({ items: [...data.items], total }, { revalidate: false });
  };

  return (
    <>
      <div className="mb-2" />
      <div>
        {data.items.map((movie, index) => (
          <MovieScreen
            key={movie.id}
            index={index}
            item={movie}
            onFixMovie={onFixMovie}
            onCopyMovie={onCopyMovie}
            onRateMovie={onRateMovie}
            onWatchedMovie={onWatchedMovie}
          />
        ))}
      </div>
    </>
  );
};
