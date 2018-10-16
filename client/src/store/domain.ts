import Vue from 'vue'
import { Module, ActionTree, MutationTree, GetterTree } from 'vuex'

import pick from 'lodash.pick'

import storage from '@/lib/storage'
import { RootState, DomainState, Movie, Movies, MovieList } from '@/types'
import * as constant from '@/constants'
import api from '@/store/api'

export const state: DomainState = {
	movies: {},
	itemsOrder: [],
	total: 0,
}

const mutations: MutationTree<DomainState> = {
	[constant.RECEIVE_MOVIES]: (local, movies: Movies) => {
		// console.log(`state-${JSON.stringify(state)}`) // eslint-disable-line
		// console.log(`total(${movies.total})-items(${movies.items.length})`) // eslint-disable-line
		local.total = movies.total // eslint-disable-line
		local.movies = movies.items.reduce(
			(list, movie: Movie) => {
				list[movie.id] = movie
				return list
			},
			{} as MovieList,
		)
		local.itemsOrder = movies.items.map((movie: Movie) => movie.id)
	},

	[constant.SET_MOVIE]: (local, movie: Movie) => {
		local.movies[movie.id] = { ...movie }
	},
}

export const actions: ActionTree<DomainState, RootState> = {
	[constant.FETCH_MOVIES_BASE]: async ({ commit, rootState }) => {
		commit(constant.SET_BUSY, true, { root: true })
		const opts = pick(rootState.options, ['query', 'filterBy', 'sortBy', 'sortOrder', 'limit', 'offset'])
		const movies: Movie[] = await api.getMovies(opts)
		commit(constant.RECEIVE_MOVIES, movies)
		commit(constant.SET_BUSY, false, { root: true })
	},

	[constant.SET_SCORE_BASE]: async (context, { id, score }) => {
		// console.log(`scoreNew(${score})-scoreOld(${context.state.movies[id].score})`) // eslint-disable-line
		const movie = { ...context.state.movies[id], score }
		// console.log(`id(${movie.id})-score(${movie.score})`) // eslint-disable-line
		context.commit(constant.SET_BUSY, true, { root: true })
		const reply: Movie = await api.setMovieScore(movie)
		context.commit(constant.SET_MOVIE, reply)
		context.commit(constant.SET_BUSY, false, { root: true })

		// api.setMovieScore(movie, changed => commit(types.SET_MOVIE, changed))
	},

	[constant.SET_WATCHED_BASE]: async (context, { id, watched }) => {
		// console.log(`scoreNew(${score})-scoreOld(${context.state.movies[id].score})`) // eslint-disable-line
		const movie = { ...context.state.movies[id], last_watched: watched }
		// console.log(`id(${movie.id})-score(${movie.score})`) // eslint-disable-line
		context.commit(constant.SET_BUSY, true, { root: true })
		const reply: Movie = await api.setMovieWatched(movie)
		context.commit(constant.SET_MOVIE, reply)
		context.commit(constant.SET_BUSY, false, { root: true })

		// api.setMovieScore(movie, changed => commit(types.SET_MOVIE, changed))
	},

	[constant.FIX_MOVIE_BASE]: async (context, { id, tmdb }) => {
		// console.log(`scoreNew(${score})-scoreOld(${context.state.movies[id].score})`) // eslint-disable-line
		const movie = { ...context.state.movies[id], tmdb_id: tmdb }
		// console.log(`id(${movie.id})-score(${movie.score})`) // eslint-disable-line
		context.commit(constant.SET_BUSY, true, { root: true })
		const reply: Movie = await api.fixMovie(movie)
		context.commit(constant.SET_MOVIE, reply)
		context.commit(constant.SET_BUSY, false, { root: true })

		// api.setMovieScore(movie, changed => commit(types.SET_MOVIE, changed))
	},

	[constant.SET_DUPLICATE_BASE]: async (context, { id, showIfDuplicate }) => {
		// console.log(`scoreNew(${score})-scoreOld(${context.state.movies[id].score})`) // eslint-disable-line
		const movie = { ...context.state.movies[id], showIfDuplicate }
		// console.log(`id(${movie.id})-score(${movie.score})`) // eslint-disable-line
		context.commit(constant.SET_BUSY, true, { root: true })
		const reply: Movie = await api.setMovieDuplicate(movie)
		context.commit(constant.SET_MOVIE, reply)
		context.commit(constant.SET_BUSY, false, { root: true })

		// api.setMovieScore(movie, changed => commit(types.SET_MOVIE, changed))
	},
}

export const getters: GetterTree<DomainState, RootState> = {
	getMovies(local: DomainState): Movie[] {
		return local.itemsOrder.map(id => local.movies[id])
	},
}

export const domain: Module<DomainState, RootState> = {
	namespaced: true,
	state,
	mutations,
	actions,
	getters,
}
