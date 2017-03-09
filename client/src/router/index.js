import Vue from 'vue';
import Router from 'vue-router';
import Cover from '@/components/Cover';

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Cover',
      component: Cover,
    },
  ],
});
