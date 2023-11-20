import React from "react";

import useSWR from "swr";
import ReactPaginate from "react-paginate";

import { getMovies, fixMovie } from "~/api";
import { useOptionsStore, useOptionsActions } from "~/state/options";
import Movie from "./movie";

const Movies = () => {
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
  );

  if (isLoading) return <div>Loading...</div>;
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
    mutate(
      { items: [...data.items], total: data.total },
      { revalidate: false },
    );
  };

  return (
    <div>
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
      <div className="mb-2" />
      <div>
        {data?.items.map((movie, index) => (
          <Movie
            key={movie.id}
            index={index}
            item={movie}
            onFixMovie={onFixMovie}
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
    </div>
  );
};

export default Movies;