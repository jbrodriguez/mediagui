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
		console.log(`total(${movies.total})-items(${movies.items.length})`) // eslint-disable-line
		local.total = movies.total // eslint-disable-line
		local.movies = movies.items.reduce(
			(list, movie: Movie) => {
				list[movie.id] = movie
				return list
			},
			{} as MovieList,
		)
		local.itemsOrder = movies.items.map((movie: Movie) => movie.id)
		// movies.items.forEach((movie: Movie) => {
		// 	state.itemsOrder.push(movie.id)
		// 	Vue.set(local.movies, movie.id, movie)
		// })
	},
	// [constant.SET_CATEGORY]: (local, category: Category) => {
	// 	// console.log(`stmt(${JSON.stringify(local)})-uno(${JSON.stringify(uno)}))-dos(${JSON.stringify(dos)})`)
	// 	// console.log(`stmt(${JSON.stringify(stmt)})`)
	// 	if (local.statement) {
	// 		local.statement.transactions[category.id].category = +category.name
	// 	}
	// },
}

export const actions: ActionTree<DomainState, RootState> = {
	[constant.FETCH_MOVIES_BASE]: async ({ commit, rootState }) => {
		commit(constant.SET_BUSY, true, { root: true })
		const opts = pick(rootState.options, ['query', 'filterBy', 'sortBy', 'sortOrder', 'limit', 'offset'])
		const movies: Movie[] = await api.getMovies(opts)
		commit(constant.RECEIVE_MOVIES, movies)
		commit(constant.SET_BUSY, false, { root: true })
	},

	// async fetchMovies({ commit, rootState }) {
	// 	commit(constant.SET_BUSY, true, { root: true })
	// 	const opts = pick(rootState.options, ['query', 'filterBy', 'sortBy', 'sortOrder', 'limit', 'offset'])
	// 	const movies: Movie[] = await api.getMovies(opts)
	// 	commit(constant.RECEIVE_MOVIES, movies)
	// 	commit(constant.SET_BUSY, false, { root: true })
	// },

	// [constant.LOAD_STATEMENT]: async (context, filepath) => {
	// 	// tslint:disable-next-line:no-console
	// 	// console.log(`file(${JSON.stringify(context)})-filepath(${filepath})`)
	// 	context.commit(constant.SET_BUSY, true)
	// 	const loaded = await api.loadStatement(filepath)
	// 	context.commit(constant.SET_STATEMENT, loaded)
	// 	context.commit(constant.SET_BUSY, false)
	// },
	// [constant.SAVE_STATEMENT]: async context => {
	// 	context.commit(constant.SET_BUSY, true)
	// 	// console.log(`statement(${JSON.stringify(context.state.statement)})`)
	// 	const loaded = await api.saveStatement(context.state.statement)
	// 	context.commit(constant.SET_BUSY, false)
	// },
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
