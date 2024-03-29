import React from "react";

import { Link } from "react-router-dom";
import debounce from "lodash.debounce";

import { useOptionsStore, useOptionsActions } from "~/state/options";
import Chevron from "~/shared/components/chevron";

export const Header: React.FC = () => {
  const { filterBy, filterByOptions, sortBy, sortByOptions } = useOptionsStore(
    (state) => ({
      filterBy: state.filterBy,
      sortBy: state.sortBy,
      filterByOptions: state.filterByOptions,
      sortByOptions: state.sortByOptions,
    }),
  );
  const { setFilterBy, setSortBy, setQuery } = useOptionsActions();

  const onFilterByChange = (e: React.ChangeEvent<HTMLSelectElement>) =>
    setFilterBy(e.target.value);

  const onSortByChange = (e: React.ChangeEvent<HTMLSelectElement>) =>
    setSortBy(e.target.value);

  const updateQuery = debounce((e: React.ChangeEvent<HTMLInputElement>) => {
    // console.log("search", e.target.value);
    setQuery(e.target.value);
  }, 750);

  return (
    <>
      <nav className="grid grid-cols-12 gap-2 py-2">
        <ul className="col-span-2 flex items-center justify-center py-2 bg-red-600 text-neutral-50">
          <li>
            <Link to="/covers">mediaGUI</Link>
          </li>
        </ul>

        <ul className="col-span-10 items-center justify-center py-2 bg-neutral-100 text-sky-700">
          <li>
            <div className="grid grid-cols-12 gap-2 justify-between">
              <div className="col-span-8 flex items-center">
                <Link to="/movies" className="mx-2">
                  MOVIES
                </Link>

                <select
                  defaultValue={filterBy}
                  className="ml-2 text-slate-600 p-1 outline-0"
                  onChange={onFilterByChange}
                >
                  {filterByOptions.map((option) => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </select>

                <input
                  type="search"
                  className="px-2 text-slate-600 border-r border-l border-slate-200 placeholder-slate-400 shadow-sm outline-0"
                  placeholder="enter search string"
                  onChange={updateQuery}
                />

                <select
                  defaultValue={sortBy}
                  className="mr-2 text-slate-600 p-1 outline-0"
                  onChange={onSortByChange}
                >
                  {sortByOptions.map((option) => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </select>

                <Chevron />

                <span className="mx-2">|</span>
                <Link to="/import">IMPORT</Link>
              </div>
              <div className="col-span-4 flex justify-end">
                <Link to="/duplicates">DUPLICATES</Link>
                <Link to="/prune" className="mx-2">
                  PRUNE
                </Link>
              </div>
            </div>
          </li>
        </ul>
      </nav>
    </>
  );
};
