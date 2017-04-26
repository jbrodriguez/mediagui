import Vue from 'vue'
import Router from 'vue-router'

import Covers from '@/components/Covers'
import Movies from '@/components/Movies'
import Import from '@/components/Import'
import Settings from '@/components/Settings'
import Duplicates from '@/components/Duplicates'
import Prune from '@/components/Prune'
import Add from '@/components/Add'

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
