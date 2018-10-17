import Vue from 'vue'
import Router from 'vue-router'

import Covers from './components/Covers.vue'
import Movies from './components/Movies.vue'
import Import from './components/Import.vue'
import Add from './components/Add.vue'
import Settings from './components/Settings.vue'
import Duplicates from './components/Duplicates.vue'
import Prune from './components/Prune.vue'

Vue.use(Router)

export default new Router({
	routes: [
		{
			path: '/',
			name: 'Covers',
			component: Covers,
		},
		{
			path: '/movies',
			name: 'Movies',
			component: Movies,
		},
		{
			path: '/import',
			name: 'Import',
			component: Import,
		},
		{
			path: '/add',
			name: 'Add',
			component: Add,
		},
		{
			path: '/settings',
			name: 'Settings',
			component: Settings,
		},
		{
			path: '/duplicates',
			name: 'Duplicates',
			component: Duplicates,
		},
		{
			path: '/prune',
			name: 'Prune',
			component: Prune,
		},
	],
})
