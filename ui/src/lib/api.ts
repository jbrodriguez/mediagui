// import React from "react";

import { type Movies } from "~/types";

type OptionsParams = {
  [key: string]: string | number | boolean;
};

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

// export async function getMovies(options: OptionsState): Promise<Movies> {
//   let movies: Movies = { total: 0, items: [] };

//   const [err, data] = await to<Response>(
//     retrieve(`${this.ep}/movies?${encode(options)}`),
//   );
//   if (err) {
//     // console.log(`reply.err(${err})`)
//   } else {
//     if (data) {
//       movies = await data.json();
//       // console.log(`data(${d})`)
//     }
//   }

//   return movies;
// }
