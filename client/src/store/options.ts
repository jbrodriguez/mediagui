import Vue from 'vue'
import { Module, ActionTree, MutationTree } from 'vuex'

import storage from '@/lib/storage'
import { RootState, OptionsState, Option } from '@/types'
import * as constant from '@/constants'
// import api from '@/store/api'
// import { Category } from '@/types'

export const state: OptionsState = {
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
	// [constant.SET_STATEMENT]: (local, stmt: Statement) => {
	// 	// console.log(`stmt(${JSON.stringify(local)})-uno(${JSON.stringify(uno)}))-dos(${JSON.stringify(dos)})`)
	// 	// console.log(`stmt(${JSON.stringify(stmt)})`)
	// 	local.statement = stmt
	// },
	// [constant.SET_CATEGORY]: (local, category: Category) => {
	// 	// console.log(`stmt(${JSON.stringify(local)})-uno(${JSON.stringify(uno)}))-dos(${JSON.stringify(dos)})`)
	// 	// console.log(`stmt(${JSON.stringify(stmt)})`)
	// 	if (local.statement) {
	// 		local.statement.transactions[category.id].category = +category.name
	// 	}
	// },
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
