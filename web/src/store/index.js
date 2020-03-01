import Vue from 'vue';
import Vuex from 'vuex';
import createPersistedState from 'vuex-persistedstate';

import board from './board';
import item from './item';
import list from './list';
import resource from './resource';
import user from './user';

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    board,
    item,
    list,
    user,
    resource,
  },
  plugins: [createPersistedState()],
});
