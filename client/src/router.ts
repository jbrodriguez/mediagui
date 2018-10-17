import Vue from 'vue'
import Router from 'vue-router'

import Covers from './components/Covers.vue'
import Movies from './components/Movies.vue'
import Import from './components/Import.vue'

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
		// {
		// 	path: '/add',
		// 	name: 'Add',
		// 	component: Add,
		// },
		// {
		// 	path: '/settings',
		// 	name: 'Settings',
		// 	component: Settings,
		// },
		// {
		// 	path: '/duplicates',
		// 	name: 'Duplicates',
		// 	component: Duplicates,
		// },
		// {
		// 	path: '/prune',
		// 	name: 'Prune',
		// 	component: Prune,
		// },
	],
	// routes: [
	// 	{
	// 		path: '/',
	// 		name: 'home',
	// 		component: Home,
	// 	},
	// 	{
	// 		path: '/about',
	// 		name: 'about',
	// 		// route level code-splitting
	// 		// this generates a separate chunk (about.[hash].js) for this route
	// 		// which is lazy-loaded when the route is visited.
	// 		component: () => import(/* webpackChunkName: "about" */ './views/About.vue'),
	// 	},
	// ],
})
