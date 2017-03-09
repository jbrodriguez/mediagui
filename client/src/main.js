// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue';
import { sync } from 'vuex-router-sync';

import 'tachyons/css/tachyons.css';
import 'flexboxgrid-sass/flexboxgrid.scss';

import './styles/styles.scss';
import App from './App';

import store from './store';
import router from './router';

sync(store, router);

Vue.config.productionTip = false;

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  template: '<App/>',
  components: { App },
});
