import Vue from 'vue'
import { Module, ActionTree, MutationTree } from 'vuex'

import storage from '@/lib/storage'
import { RootState, OptionsState, Option } from '@/types'
import * as constant from '@/constants'

export const state: OptionsState = {
	query: '',
	filterByOptions: [
		{ value: 'title', label: 'Title' },
		{ value: 'genre', label: 'Genre' },
		{ value: 'year', label: 'Year' },
		{ value: 'country', label: 'Country' },
		{ value: 'director', label: 'Director' },
		{ value: 'actor', label: 'Actor' },
		{ value: 'location', label: 'Location' },
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
		{ value: 'score', label: 'Score' },
	],
	sortBy: storage.get('sortBy') || 'added',
	sortOrderOptions: ['asc', 'desc'],
	sortOrder: storage.get('sortOrder') || 'desc',
	mode: 'regular',
	limit: 50,
	offset: 0,
}

const mutations: MutationTree<OptionsState> = {
	[constant.SET_FILTER]: (local, filterBy) => {
		local.filterBy = filterBy
		local.offset = 0
		storage.set('filterBy', filterBy)
	},

	[constant.SET_QUERY]: (local, query) => {
		local.query = query
		local.offset = 0
	},

	[constant.SET_SORT]: (local, sortBy) => {
		local.sortBy = sortBy
		storage.set('sortBy', sortBy)
	},

	[constant.SET_OFFSET]: (local, offset) => {
		local.offset = offset
	},

	[constant.FLIP_ORDER]: (local, sortOrder) => {
		local.sortOrder = local.sortOrder === 'asc' ? 'desc' : 'asc'
		storage.set('sortOrder', local.sortOrder)
	},
}

export const actions: ActionTree<OptionsState, RootState> = {
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

export const options: Module<OptionsState, RootState> = {
	state,
	mutations,
	actions,
}
