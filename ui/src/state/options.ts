import { createWithEqualityFn } from "zustand/traditional";
import { persist } from "zustand/middleware";
import { shallow } from "zustand/shallow";

interface OptionsState {
  query: string;
  filterByOptions: { value: string; label: string }[];
  filterBy: string;
  sortByOptions: { value: string; label: string }[];
  sortBy: string;
  sortOrderOptions: string[];
  sortOrder: string;
  mode: string;
  limit: number;
  offset: number;
  actions: {
    setFilterBy: (filterBy: string) => void;
    setQuery: (query: string) => void;
    setSortBy: (sortBy: string) => void;
    setOffset: (offset: number) => void;
    changeOffset: (index: number) => void;
    flipOrder: () => void;
  };
}

export const useOptionsStore = createWithEqualityFn<OptionsState>()(
  persist(
    (set) => ({
      query: "",
      filterByOptions: [
        { value: "title", label: "Title" },
        { value: "genre", label: "Genre" },
        { value: "year", label: "Year" },
        { value: "country", label: "Country" },
        { value: "director", label: "Director" },
        { value: "actor", label: "Actor" },
        { value: "location", label: "Location" },
      ],
      filterBy: "title",
      sortByOptions: [
        { value: "title", label: "Title" },
        { value: "runtime", label: "Runtime" },
        { value: "added", label: "Added" },
        { value: "last_watched", label: "Watched W" },
        { value: "count_watched", label: "Watched C" },
        { value: "year", label: "Year" },
        { value: "imdb_rating", label: "Rating" },
        { value: "score", label: "Score" },
      ],
      sortBy: "added",
      sortOrderOptions: ["asc", "desc"],
      sortOrder: "desc",
      mode: "regular",
      limit: 50,
      offset: 0,
      actions: {
        setQuery: (query) => set({ query, offset: 0 }),
        setFilterBy: (filterBy) => set({ filterBy, offset: 0 }),
        setSortBy: (sortBy) => set({ sortBy }),
        setOffset: (offset) => set({ offset }),
        changeOffset: (index) =>
          set((state) => ({ offset: Math.ceil(index * state.limit) })),
        flipOrder: () =>
          set((state) => ({
            sortOrder: state.sortOrder === "asc" ? "desc" : "asc",
          })),
      },
    }),
    {
      name: "options",
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      partialize: ({
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        query,
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        filterByOptions,
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        sortByOptions,
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        sortOrderOptions,
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        limit,
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        offset,
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        actions,
        ...rest
      }) => rest,
    },
  ),
  shallow,
);

export const useOptionsActions = () =>
  useOptionsStore((state) => state.actions);

export const useFilters = () =>
  useOptionsStore((state) => state.filterByOptions);

export const useOptionsSortOrder = () =>
  useOptionsStore((state) => state.sortOrder);
