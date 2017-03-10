import Vue from 'vue';
import Router from 'vue-router';
import Covers from '@/components/Covers';

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Covers',
      component: Covers,
    },
  ],
});
