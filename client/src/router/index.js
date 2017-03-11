import Vue from 'vue';
import Router from 'vue-router';
import Covers from '@/components/Covers';
import Movies from '@/components/Movies';

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
  ],
});
