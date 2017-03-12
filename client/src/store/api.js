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
      body: JSON.stringify({ topic: '', payload: folder }),
    })
    .then(resp => resp.json())
    .then(data => cb(data));
  }

  pruneMovies() {
    return fetch(`${this.host}/prune`, { method: 'POST' });
  }
}

const api = new Api();

export default api;
