import Vue from 'vue';
import Vuex from 'vuex';

import api from './api';
import * as types from './types';

Vue.use(Vuex);

const store = new Vuex.Store({
  state: {
    movies: { /* [id: number] Movie */ },
  },

  actions: {
    [types.FETCH_COVER]: ({ commit }) =>
      api.getCover((movies) => {
        // console.log(`movies-${JSON.stringify(movies)}`); // eslint-disable-line
        commit(types.RECEIVE_MOVIES, movies);
      }),
  },

  mutations: {
    [types.RECEIVE_MOVIES]: (state, movies) => {
      console.log(`state-${JSON.stringify(state)}`); // eslint-disable-line
      console.log(`context-${movies.total}`);  // eslint-disable-line
      movies.items.forEach((movie) => {
        Vue.set(state.movies, movie.id, movie);
      });
    },
  },

  getters: {
    getMovies(state) {
      console.log(`length-${Object.keys(state.movies).length}`); // eslint-disable-line
      return Object.keys(state.movies).map(id => state.movies[id]).reverse();
    },
  },
});

export default store;
