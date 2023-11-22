import React from "react";

import useSWR from "swr";
import ReactPaginate from "react-paginate";

import { getMovies, fixMovie, copyMovie, rateMovie, watchedMovie } from "~/api";
import { useOptionsStore, useOptionsActions } from "~/state/options";
import { MovieScreen } from "~/shared/movie/movie";
import { Spinner } from "~/shared/components/spinner";

export const Movies = () => {
  const [pageIndex, setPageIndex] = React.useState(0);

  const { query, filterBy, sortBy, sortOrder, limit, offset } =
    useOptionsStore();
  const { changeOffset } = useOptionsActions();

  const { data, error, mutate, isLoading } = useSWR(
    {
      url: "/movies",
      args: { query, filterBy, sortBy, sortOrder, limit, offset },
    },
    getMovies,
    {
      keepPreviousData: true,
    },
  );

  if (error) return <div>Error</div>;
  if (!data) return <div>No data</div>;

  // console.log("data", data);
  const total = data.total ?? 0;
  const pageCount = Math.ceil(total / 50);

  const handlePageClick = (e: { selected: number }) => {
    changeOffset(e.selected);
    setPageIndex(e.selected);
  };

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
    mutate({ items: data.items, total }, { revalidate: false });
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
    mutate({ items: data.items, total }, { revalidate: false });
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
    mutate({ items: data.items, total }, { revalidate: false });
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
    mutate({ items: data.items, total }, { revalidate: false });
  };

  return (
    <>
      <div className="flex items-center justify-between">
        <ReactPaginate
          breakLabel="..."
          nextLabel="Next"
          onPageChange={handlePageClick}
          pageRangeDisplayed={5}
          pageCount={pageCount}
          previousLabel="Prev"
          renderOnZeroPageCount={null}
          forcePage={pageIndex}
          disableInitialCallback={true}
          containerClassName="flex flex-row justify-start items-center"
          pageClassName="px-1"
          pageLinkClassName="px-4 py-1 flex items-center justify-center p-0 text-gray-500 transition duration-150 ease-in-out hover:bg-light-300"
          activeLinkClassName="border bg-sky-600 text-neutral-100 cursor-default"
          breakLinkClassName="text-gray-500"
          previousLinkClassName="pr-4 py-1 flex items-center justify-center p-0 text-gray-500 transition duration-150 ease-in-out hover:bg-light-300"
          nextLinkClassName="px-4 py-1 flex items-center justify-center p-0 text-gray-500 transition duration-150 ease-in-out hover:bg-light-300"
        />
        <div className="flex flex-row items-center">
          {isLoading ? <Spinner /> : null}
          <span className="text-slate-500">TOTAL</span>
          <span className="text-sky-950 font-bold ml-1">{total}</span>
        </div>
      </div>
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
      <ReactPaginate
        breakLabel="..."
        nextLabel="Next"
        onPageChange={handlePageClick}
        pageRangeDisplayed={5}
        pageCount={pageCount}
        previousLabel="Prev"
        renderOnZeroPageCount={null}
        forcePage={pageIndex}
        disableInitialCallback={true}
        containerClassName="flex flex-row justify-start items-center"
        pageClassName="px-1"
        pageLinkClassName="px-4 py-1 flex items-center justify-center p-0 text-gray-500 transition duration-150 ease-in-out hover:bg-light-300"
        activeLinkClassName="border bg-sky-600 text-neutral-100 cursor-default"
        breakLinkClassName="text-gray-500"
        previousLinkClassName="pr-4 py-1 flex items-center justify-center p-0 text-gray-500 transition duration-150 ease-in-out hover:bg-light-300"
        nextLinkClassName="px-4 py-1 flex items-center justify-center p-0 text-gray-500 transition duration-150 ease-in-out hover:bg-light-300"
      />
    </>
  );
};
