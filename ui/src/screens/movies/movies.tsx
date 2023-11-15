// import React from "react";

import useSWR from "swr";

import { getMovies } from "~/lib/api";
import { useOptionsStore } from "~/state/options";
import type { Movies, Movie } from "~/types";

const Movies = () => {
  // const [pageIndex, _] = React.useState(0);
  const { query, filterBy, sortBy, sortOrder, limit, offset } =
    useOptionsStore();
  const args = { query, filterBy, sortBy, sortOrder, limit, offset };

  console.log("movies params ", args);

  const { data, error, isLoading } = useSWR(
    { url: "/movies", args },
    getMovies,
  );

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error</div>;

  console.log("movies data total", data);

  return (
    <div>
      <h1>Movies</h1>
      <h2>{data?.total}</h2>
      <div>
        {data?.items.map((movie: Movie) => (
          <div key={movie.id}>
            <h3>{movie.title}</h3>
            <p>{movie.score}</p>
          </div>
        ))}
      </div>
      {/* Add your movie content here */}
    </div>
  );
};

export default Movies;
