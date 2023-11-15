// workaround - https://forum.vuejs.org/t/vuex-submodules-with-typescript/40903/4
export interface RootState {
  isBusy: boolean;
  lines: string[];
  config: ConfigState | object;
  options: OptionsState | object;
  domain: DomainState | object;
}

export interface ConfigState {
  version: string;
  unraidMode: boolean;
  unraidHosts: string[];
  mediaFolders: string[];
}

export interface OptionsState {
  query: string;
  filterByOptions: Option[];
  filterBy: string;
  sortByOptions: Option[];
  sortBy: string;
  sortOrderOptions: ["asc", "desc"];
  sortOrder: string;
  mode: string;
  limit: number;
  offset: number;
}

export interface OptionsParams {
  query: string;
  filterBy: string;
  sortBy: string;
  sortOrder: string;
  limit: number;
  offset: number;
}

export interface DomainState {
  movies: MovieList;
  itemsOrder: number[];
  total: number;
}

export interface Option {
  value: string;
  label: string;
}

export interface MovieList {
  [id: number]: Movie;
}

export interface Movie {
  id: number;
  title: string;
  runtime: number;
  count_watched: number;
  tmdb_id: number;
  score: number;
  added: string;
  last_watched: string;
  cover: string;
  backdrop: string;
  overview: string;
  all_watched: string;
  showIfDuplicate: number;
}

export interface Movies {
  total: number;
  items: Movie[];
}
