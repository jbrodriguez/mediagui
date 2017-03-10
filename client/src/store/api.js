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
}

const api = new Api();

export default api;
