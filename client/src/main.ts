import Vue from 'vue'

import '@/styles/styles.scss'

import { library } from '@fortawesome/fontawesome-svg-core'
import {
	faChevronCircleUp,
	faChevronCircleDown,
	faStar,
	faBinoculars,
	faPlus,
	faTimesCircle,
} from '@fortawesome/free-solid-svg-icons'
import { faStar as faStarO } from '@fortawesome/free-regular-svg-icons'

// https://github.com/FortAwesome/vue-fontawesome/issues/24#issuecomment-417897681
/* tslint:disable:no-var-requires */
const fontawesome = require('@fortawesome/vue-fontawesome')
// import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

import App from './App.vue'
import router from './router'
import store from './store'

library.add(faChevronCircleUp, faChevronCircleDown, faStar, faBinoculars, faPlus, faTimesCircle, faStarO)

Vue.component('font-awesome-icon', fontawesome.FontAwesomeIcon)

Vue.config.productionTip = false

new Vue({
	router,
	store,
	render: h => h(App),
}).$mount('#app')
