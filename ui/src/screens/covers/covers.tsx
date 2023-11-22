import useSWR from "swr";

import { getMovies } from "~/api";
import { useOptionsStore } from "~/state/options";
import { CoverScreen } from "./cover";
// import { Spinner } from "~/shared/components/spinner";

export const CoversScreen = () => {
  const { query, filterBy, sortBy, sortOrder } = useOptionsStore();

  const { data, error } = useSWR(
    {
      url: "/movies",
      args: { query, filterBy, sortBy, sortOrder, limit: 60, offset: 0 },
    },
    getMovies,
    {
      keepPreviousData: true,
    },
  );

  if (error) return <div>Error</div>;
  if (!data) return <div>No data</div>;

  // console.log("data", data);
  // const total = data.total ?? 0;

  return (
    <>
      <div className="mb-2" />
      <div className="grid grid-cols-6 gap-2">
        {data.items.map((movie) => (
          <CoverScreen key={movie.id} movie={movie} />
        ))}
      </div>
    </>
  );
};
