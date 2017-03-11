import Vue from 'vue';
import Vuex from 'vuex';

import pick from 'lodash.pick';

import api from './api';
import * as types from './types';
import storage from '../lib/storage';

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
  },

  mutations: {
    [types.RECEIVE_CONFIG]: (state, config) => {
      state.config = { ...config }; // eslint-disable-line
    },

    [types.RECEIVE_MOVIES]: (state, movies) => {
      // console.log(`state-${JSON.stringify(state)}`); // eslint-disable-line
      // console.log(`context-${movies.total}`);  // eslint-disable-line
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
});

export default store;
