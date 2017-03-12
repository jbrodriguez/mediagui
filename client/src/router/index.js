import Vue from 'vue';
import Router from 'vue-router';

import Covers from '@/components/Covers';
import Movies from '@/components/Movies';
import Import from '@/components/Import';

Vue.use(Router);

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
  ],
});
