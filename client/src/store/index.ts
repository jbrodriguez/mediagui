import Vue from 'vue'
import Vuex, { StoreOptions, ActionTree, MutationTree } from 'vuex'

import { RootState, Movie } from '@/types'
import * as constant from '@/constants'
import api from '@/store/api'
import socket from '@/store/socket'

import { options } from './options'
import { config } from './config'
import { domain } from './domain'

Vue.use(Vuex)

const socketPlugin = local => {
	socket.receive(message => {
		const packet = JSON.parse(message.data)
		if (typeof packet.topic === 'string' && packet.topic.length > 0) {
			local.commit(packet.topic, packet.payload)
		}
	})

	local.subscribe(mutation => {
		if (mutation.type.startsWith('skt')) {
			socket.send({ topic: mutation.type, payload: mutation.payload })
		}
	})
}

const mutations: MutationTree<RootState> = {
	[constant.SET_BUSY]: (state, isBusy: boolean) => {
		state.isBusy = isBusy
	},

	[constant.IMPORT_BEGIN]: (state, line) => {
		state.lines = [line]
	},

	[constant.IMPORT_PROGRESS]: (state, line) => {
		state.lines.push(line)
	},

	[constant.IMPORT_END]: (state, line) => {
		state.lines.push(line)
	},

	[constant.PRUNE_BEGIN]: (state, line) => {
		state.lines = [line]
	},

	[constant.PRUNE_SELECTED]: (state, line) => {
		state.lines.push(line)
	},

	[constant.PRUNE_DELETE]: (state, line) => {
		state.lines.push(line)
	},

	[constant.PRUNE_END]: (state, line) => {
		state.lines.push(line)
	},
}

const actions: ActionTree<RootState, RootState> = {
	[constant.RUN_IMPORT]: context => {
		api.importMovies()
	},

	[constant.RUN_PRUNE]: context => {
		api.pruneMovies()
	},
}

const store: StoreOptions<RootState> = {
	state: {
		isBusy: false,
		lines: [],
		options: {},
		config: {},
		domain: {},
	},
	mutations,
	actions,
	modules: {
		options,
		config,
		domain,
	},
	plugins: [socketPlugin],
}

export default new Vuex.Store<RootState>(store)
