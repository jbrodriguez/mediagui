const encode = (data) => {
  const encoded = Object.keys(data).map((key) => {
    const value = encodeURIComponent(data[key].toString());
    return `${key}=${value}`;
  });
  return encoded.join('&');
};

class Api {
  host = `http://${document.location.host}/api/v1`;
  // host = 'http://blackbeard.apertoire.org:7623/api/v1';

  getConfig(cb) {
    return fetch(`${this.host}/config`)
          .then(resp => resp.json())
          .then(data => cb(data));
  }

  getCover(cb) {
    return fetch(`${this.host}/movies/cover`)
          .then(resp => resp.json())
          .then(data => cb(data));
  }

  getMovies(options, cb) {
    return fetch(`${this.host}/movies?${encode(options)}`)
          .then(resp => resp.json())
          .then(data => cb(data));
  }

  importMovies() {
    return fetch(`${this.host}/import`, { method: 'POST' });
  }

  addFolder(folder, cb) {
    return fetch(`${this.host}/config/folder`, {
      method: 'PUT',
      headers: new Headers({ 'Content-Type': 'application/json' }),
      body: JSON.stringify({ topic: '', payload: folder }),
    })
    .then(resp => resp.json())
    .then(data => cb(data));
  }

  pruneMovies() {
    return fetch(`${this.host}/prune`, { method: 'POST' });
  }

  getDuplicates(cb) {
    return fetch(`${this.host}/movies/duplicates`)
          .then(resp => resp.json())
          .then(data => cb(data));
  }

  setMovieScore(movie, cb) {
    return fetch(`${this.host}/movies/${movie.id}/score`, {
      method: 'PUT',
      headers: new Headers({ 'Content-Type': 'application/json' }),
      body: JSON.stringify(movie),
    })
    .then(resp => resp.json())
    .then(data => cb(data));
  }

  setMovieWatched(movie, cb) {
    return fetch(`${this.host}/movies/${movie.id}/watched`, {
      method: 'PUT',
      headers: new Headers({ 'Content-Type': 'application/json' }),
      body: JSON.stringify(movie),
    })
    .then(resp => resp.json())
    .then(data => cb(data));
  }

  fixMovie(movie, cb) {
    return fetch(`${this.host}/movies/${movie.id}/fix`, {
      method: 'PUT',
      headers: new Headers({ 'Content-Type': 'application/json' }),
      body: JSON.stringify(movie),
    })
    .then(resp => resp.json())
    .then(data => cb(data));
  }

  setMovieDuplicate(movie, cb) {
    return fetch(`${this.host}/movies/${movie.id}/duplicate`, {
      method: 'PUT',
      headers: new Headers({ 'Content-Type': 'application/json' }),
      body: JSON.stringify(movie),
    })
    .then(resp => resp.json())
    .then(data => cb(data));
  }
}

const api = new Api();

export default api;
