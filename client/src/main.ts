import Vue from 'vue'

import 'tachyons/css/tachyons.min.css'
import 'flexboxgrid-sass/flexboxgrid.scss'

import { library } from '@fortawesome/fontawesome-svg-core'
import { faChevronCircleUp, faChevronCircleDown } from '@fortawesome/free-solid-svg-icons'

// https://github.com/FortAwesome/vue-fontawesome/issues/24#issuecomment-417897681
/* tslint:disable:no-var-requires */
const fontawesome = require('@fortawesome/vue-fontawesome')
// import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

import App from './App.vue'
import router from './router'
import store from './store'

library.add(faChevronCircleUp, faChevronCircleDown)

Vue.component('font-awesome-icon', fontawesome.FontAwesomeIcon)

Vue.config.productionTip = false

new Vue({
	router,
	store,
	render: h => h(App),
}).$mount('#app')
