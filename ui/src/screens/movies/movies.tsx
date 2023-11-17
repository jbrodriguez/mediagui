import React from "react";

import useSWR from "swr";
import ReactPaginate from "react-paginate";

import { getMovies } from "~/api";
import { useOptionsStore, useOptionsActions } from "~/state/options";
import Movie from "./movie";

const Movies = () => {
  const [pageIndex, setPageIndex] = React.useState(0);

  const { query, filterBy, sortBy, sortOrder, limit, offset } =
    useOptionsStore();
  const { changeOffset } = useOptionsActions();

  const { data, error, isLoading } = useSWR(
    {
      url: "/movies",
      args: { query, filterBy, sortBy, sortOrder, limit, offset },
    },
    getMovies,
  );

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error</div>;

  // console.log("data", data);

  const total = data?.total ?? 0;
  const pageCount = Math.ceil(total / 50);

  const handlePageClick = (e: { selected: number }) => {
    changeOffset(e.selected);
    setPageIndex(e.selected);
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
        {data?.items.map((movie) => <Movie key={movie.id} item={movie} />)}
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
