import Vue from 'vue';
import VueRouter from 'vue-router';
import PageBoards from '../views/PageBoards.vue';
import PageBoard from '../views/PageBoard.vue';

Vue.use(VueRouter);

const routes = [
  {
    path: '/',
    name: 'boards',
    component: PageBoards,
  },
  {
    path: '/board/:id',
    name: 'board',
    component: PageBoard,
  },
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
});

export default router;
