import Vue from 'vue';
import VueRouter from 'vue-router';
import BoardsPage from '../views/BoardsPage.vue';
import BoardPage from '../views/BoardPage.vue';

Vue.use(VueRouter);

const routes = [
  {
    path: '/',
    name: 'boards',
    component: BoardsPage,
  },
  {
    path: '/board/:id',
    name: 'board',
    component: BoardPage,
  },
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
});

export default router;
