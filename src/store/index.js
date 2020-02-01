import Vue from 'vue';
import Vuex from 'vuex';

import board from './board';
import item from './item';
import list from './list';
import user from './user';

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    board,
    item,
    list,
    user,
  },
});
