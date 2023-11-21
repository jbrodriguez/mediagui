// import React from "react";

import type { Movies, OptionsParams, Movie } from "~/types";

const encode = (params: OptionsParams): string => {
  const str = [];
  for (const key in params) {
    if (Object.prototype.hasOwnProperty.call(params, key)) {
      str.push(`${encodeURIComponent(key)}=${encodeURIComponent(params[key])}`);
    }
  }
  return str.join("&");
};

const apiEndpoint = `${document.location.origin}/api/v1`;

export async function getMovies(params: {
  url: string;
  args: OptionsParams;
}): Promise<Movies> {
  console.log("getMovies", params.url, params.args);
  console.log("document.location", document.location);
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
