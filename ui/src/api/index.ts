// import React from "react";

import type { Movies, OptionsParams, Movie, ConfigState } from "~/types";

const encode = (params: OptionsParams): string => {
  const str = [];
  for (const key in params) {
    if (Object.prototype.hasOwnProperty.call(params, key)) {
      str.push(`${encodeURIComponent(key)}=${encodeURIComponent(params[key])}`);
    }
  }
  return str.join("&");
};

export const apiEndpoint = `${document.location.origin}/api/v1`;

export const wsEndpoint = `${
  document.location.protocol === "http:" ? "ws:" : "wss:"
}//${document.location.host}/ws`;

export async function getConfig(): Promise<ConfigState> {
  const response = await fetch(`${apiEndpoint}/config`);
  // if (!response.ok) {
  //   throw new Error(response.statusText);
  // }

  return await response.json();
}

export async function getMovies(params: {
  url: string;
  args: OptionsParams;
}): Promise<Movies> {
  const response = await fetch(
    `${apiEndpoint}${params.url}?${encode(params.args)}`,
  );
  // if (!response.ok) {
  //   throw new Error(response.statusText);
  // }

  return await response.json();
}

export async function fixMovie(params: {
  id: number;
  tmdb_id: number;
}): Promise<Movie> {
  const response = await fetch(`${apiEndpoint}/movies/${params.id}/fix`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(params),
  });
  // if (!response.ok) {
  //   throw new Error(response.statusText);
  // }

  return await response.json();
}

export async function copyMovie(params: {
  id: number;
  tmdb_id: number;
}): Promise<Movie> {
  const response = await fetch(`${apiEndpoint}/movies/${params.id}/copy`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(params),
  });
  // if (!response.ok) {
  //   throw new Error(response.statusText);
  // }

  return await response.json();
}

export async function rateMovie(params: {
  id: number;
  score: number; // score
}): Promise<Movie> {
  const response = await fetch(`${apiEndpoint}/movies/${params.id}/score`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(params),
  });
  // if (!response.ok) {
  //   throw new Error(response.statusText);
  // }

  return await response.json();
}

export async function watchedMovie(params: {
  id: number;
  watched: string;
}): Promise<Movie> {
  const response = await fetch(`${apiEndpoint}/movies/${params.id}/watched`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(params),
  });
  // if (!response.ok) {
  //   throw new Error(response.statusText);
  // }

  return await response.json();
}

export async function importMovies() {
  fetch(`${apiEndpoint}/import`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  });
  // if (!response.ok) {
  //   throw new Error(response.statusText);
  // }

  return;
}

export async function pruneMovies() {
  fetch(`${apiEndpoint}/prune`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  });
  // if (!response.ok) {
  //   throw new Error(response.statusText);
  // }

  return;
}

export async function getDuplicates(url: string): Promise<Movies> {
  const response = await fetch(`${apiEndpoint}${url}`);
  // if (!response.ok) {
  //   throw new Error(response.statusText);
  // }

  return await response.json();
}
