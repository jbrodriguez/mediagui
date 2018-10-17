import Vue from 'vue'
import { Module, ActionTree, MutationTree } from 'vuex'

import storage from '@/lib/storage'
import { RootState, ConfigState } from '@/types'
import * as constant from '@/constants'
import api from '@/store/api'

export const state: ConfigState = {
	version: 'loading...',
	unraidMode: true,
	unraidHosts: [],
	mediaFolders: [],
}

const mutations: MutationTree<ConfigState> = {
	[constant.RECEIVE_CONFIG]: (local, conf: ConfigState) => {
		// console.log(`local(${JSON.stringify(local)})-conf(${JSON.stringify(conf)}))})`)
		// console.log(`stmt(${JSON.stringify(stmt)})`)
		// local = conf
		local.version = conf.version
		local.unraidMode = conf.unraidMode
		local.unraidHosts = conf.unraidHosts
		local.mediaFolders = conf.mediaFolders
	},
	// [constant.SET_CATEGORY]: (local, category: Category) => {
	// 	// console.log(`stmt(${JSON.stringify(local)})-uno(${JSON.stringify(uno)}))-dos(${JSON.stringify(dos)})`)
	// 	// console.log(`stmt(${JSON.stringify(stmt)})`)
	// 	if (local.statement) {
	// 		local.statement.transactions[category.id].category = +category.name
	// 	}
	// },
}

export const actions: ActionTree<ConfigState, RootState> = {
	[constant.FETCH_CONFIG]: async context => {
		context.commit(constant.SET_BUSY, true)
		const conf: ConfigState = await api.getConfig()
		context.commit(constant.RECEIVE_CONFIG, conf)
		context.commit(constant.SET_BUSY, false)
		// // tslint:disable-next-line:no-console
		// console.log(`config(${JSON.stringify(conf)})`)
	},

	[constant.ADD_FOLDER]: async (context, folder: string) => api.addFolder(folder),

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

export const config: Module<ConfigState, RootState> = {
	state,
	mutations,
	actions,
}
