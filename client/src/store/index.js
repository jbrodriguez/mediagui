import Vue from 'vue';
import Vuex from 'vuex';

import pick from 'lodash.pick';

import api from './api';
import socket from './socket';
import * as types from './types';
import storage from '../lib/storage';

// websocket handler
const socketPlugin = (store) => {
  socket.receive((message) => {
    const packet = JSON.parse(message.data);
    if (typeof packet.topic === 'string' && packet.topic.length > 0) {
      store.commit(packet.topic, packet.payload);
    }
  });

  store.subscribe((mutation) => {
    if (mutation.type.startsWith('skt')) {
      socket.send({ topic: mutation.type, payload: mutation.payload });
    }
  });
};

Vue.use(Vuex);

const store = new Vuex.Store({
  state: {
    movies: { /* [id: number] Movie */ },
    config: {},
    options: {
      query: '',
      filterByOptions: [
        { value: 'title', label: 'Title' },
        { value: 'genre', label: 'Genre' },
        { value: 'year', label: 'Year' },
        { value: 'country', label: 'Country' },
        { value: 'director', label: 'Director' },
        { value: 'actor', label: 'Actor' },
      ],
      filterBy: storage.get('filterBy') || 'title',
      sortByOptions: [
        { value: 'title', label: 'Title' },
        { value: 'runtime', label: 'Runtime' },
        { value: 'added', label: 'Added' },
        { value: 'last_watched', label: 'Watched W' },
        { value: 'count_watched', label: 'Watched C' },
        { value: 'year', label: 'Year' },
        { value: 'imdb_rating', label: 'Rating' },
      ],
      sortBy: storage.get('sortBy') || 'added',
      sortOrderOptions: ['asc', 'desc'],
      sortOrder: storage.get('sortOrder') || 'desc',
      mode: 'regular',
      limit: 50,
      offset: 0,
    },
    itemsOrder: [],
    lines: [],
  },

  actions: {
    [types.FETCH_CONFIG]: ({ commit }) =>
      api.getConfig((config) => {
        commit(types.RECEIVE_CONFIG, config);
      }),

    [types.FETCH_COVER]: ({ commit }) =>
      api.getCover((movies) => {
        // console.log(`movies-${JSON.stringify(movies)}`); // eslint-disable-line
        commit(types.RECEIVE_MOVIES, movies);
      }),

    [types.FETCH_MOVIES]: ({ commit, state }) => {
      const opts = pick(state.options, [
        'query', 'filterBy', 'sortBy', 'sortOrder', 'limit', 'offset',
      ]);

      api.getMovies(opts, (movies) => {
        // console.log(`movies-${JSON.stringify(movies)}`); // eslint-disable-line
        commit(types.RECEIVE_MOVIES, movies);
      });
    },

    [types.RUN_IMPORT]: () => api.importMovies(),

    [types.ADD_FOLDER]: folder => api.addFolder(folder),

    [types.RUN_PRUNE]: () => api.pruneMovies(),

    [types.FETCH_DUPLICATES]: ({ commit }) => {
      api.getDuplicates(movies => commit(types.RECEIVE_MOVIES, movies));
    },

    [types.SET_SCORE]: ({ commit, state }, { id, score }) => {
      console.log(`score(${score})`); // eslint-disable-line
      const movie = Object.assign({}, state.movies[id], { score });
      console.log(`id(${movie.id})-score(${movie.score})`); // eslint-disable-line
      api.setMovieScore(movie, changed => commit(types.SET_MOVIE, changed));
    },

    [types.SET_WATCHED]: ({ commit, state }, { id, watched }) => {
      const movie = Object.assign({}, state.movies[id], { last_watched: watched });
      console.log(`id(${movie.id})-last_watched(${movie.last_watched})`); // eslint-disable-line
      api.setMovieWatched(movie, changed => commit(types.SET_MOVIE, changed));
    },

    [types.FIX_MOVIE]: ({ commit, state }, { id, tmdb }) => {
      const movie = Object.assign({}, state.movies[id], { tmdb_id: tmdb });
      console.log(`id(${movie.id})-last_watched(${movie.tmdb_id})`); // eslint-disable-line
      api.fixMovie(movie, changed => commit(types.SET_MOVIE, changed));
    },
  },

  mutations: {
    [types.RECEIVE_CONFIG]: (state, config) => {
      state.config = { ...config }; // eslint-disable-line
    },

    [types.RECEIVE_MOVIES]: (state, movies) => {
      // console.log(`state-${JSON.stringify(state)}`); // eslint-disable-line
      console.log(`total(${movies.total})-items(${movies.items.length})`);  // eslint-disable-line
      state.itemsOrder = []; // eslint-disable-line
      movies.items.forEach((movie) => {
        state.itemsOrder.push(movie.id);
        Vue.set(state.movies, movie.id, movie);
      });
    },

    [types.SET_FILTER]: (state, filterBy) => {
      state.options.filterBy = filterBy; // eslint-disable-line
      storage.set('filterBy', filterBy);
    },

    [types.SET_QUERY]: (state, query) => {
      state.options.query = query; // eslint-disable-line
    },

    [types.SET_SORT]: (state, sortBy) => {
      state.options.sortBy = sortBy; // eslint-disable-line
      storage.set('sortBy', sortBy);
    },

    [types.FLIP_ORDER]: (state) => {
      state.options.sortOrder = state.options.sortOrder === 'asc' ? 'desc' : 'asc'; // eslint-disable-line
      storage.set('sortOrder', state.options.sortOrder);
    },

    [types.SET_MOVIE]: (state, movie) => {
      Vue.set(state.movies, movie.id, movie);
    },

    [types.IMPORT_BEGIN]: (state, line) => {
      state.lines = [line]; // eslint-disable-line
    },

    [types.IMPORT_PROGRESS]: (state, line) => {
      state.lines.push(line);
    },

    [types.IMPORT_END]: (state, line) => {
      state.lines.push(line);
    },

    [types.PRUNE_BEGIN]: (state, line) => {
      state.lines = [line]; // eslint-disable-line
    },

    [types.PRUNE_SELECTED]: (state, line) => {
      state.lines.push(line);
    },

    [types.PRUNE_DELETE]: (state, line) => {
      state.lines.push(line);
    },

    [types.PRUNE_END]: (state, line) => {
      state.lines.push(line);
    },
  },

  getters: {
    getMovies(state) {
      // console.log(`length-${Object.keys(state.movies).length}`); // eslint-disable-line
      // return Object.keys(state.movies).map(id => state.movies[id]).reverse();
      return state.itemsOrder.map(id => state.movies[id]);
    },

    version(state) {
      return state.config ? state.config.version : '';
    },
  },

  plugins: [socketPlugin],
});

export default store;
