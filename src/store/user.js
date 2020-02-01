// interface User {
//   id: '',
//   name: ''
// }

function state() {
  return {
    user: {},
  };
}

const mutations = {
  rename({ user }, newName) {
    user.name = newName;
  },
};

const actions = {
  rename({ commit }, newName) {
    commit('rename', newName);
  },
};

const getters = {
  getUser({ user }) {
    return user;
  },
};

export default {
  state,
  mutations,
  actions,
  getters,
};
