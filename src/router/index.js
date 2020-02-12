import Vue from 'vue';
import VueRouter from 'vue-router';
import PageBoards from '../views/PageBoards.vue';
import PageBoard from '../views/PageBoard.vue';
import PageHome from '../views/PageHome.vue';

Vue.use(VueRouter);

const routes = [
  {
    path: '/boards',
    name: 'boards',
    component: PageBoards,
  },
  {
    path: '/',
    name: 'home',
    component: PageHome,
  },
  {
    path: '/boards/:id',
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
