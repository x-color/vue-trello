// interface User {
//   id: string;
//   name: string;
//   boards: []string;
// }

function state() {
  return {
    user: {
      id: '0',
      name: 'user1',
      login: true,
      boards: [],
    },
  };
}

const mutations = {
  editUser(_state, newUser) {
    _state.user = newUser;
  },
};

const actions = {
  editUser({ commit }, newUser) {
    commit('editUser', newUser);
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
