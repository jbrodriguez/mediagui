import Vue from 'vue'
import Vuex, { StoreOptions, ActionTree, MutationTree } from 'vuex'

import { RootState, Movie } from '@/types'
import * as constant from '@/constants'
import api from '@/store/api'

import { options } from './options'
import { config } from './config'
import { domain } from './domain'

Vue.use(Vuex)

const mutations: MutationTree<RootState> = {
	[constant.SET_BUSY]: (state, isBusy: boolean) => {
		state.isBusy = isBusy
	},
}

const store: StoreOptions<RootState> = {
	state: {
		isBusy: false,
		options: {},
		config: {},
		domain: {},
	},
	mutations,
	modules: {
		options,
		config,
		domain,
	},
}

export default new Vuex.Store<RootState>(store)
