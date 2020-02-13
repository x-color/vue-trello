// interface User {
//   id: string;
//   name: string;
//   boards: []string;
// }

function state() {
  return {
    user: {
      id: '0',
      name: 'test user',
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
  login({ commit, getters }, { username, password }) {
    // Dummy login request
    if (!(username === 'testuser' && password === 'password')) {
      return false;
    }
    // Dummy logged in user data
    const loggedInUser = { ...getters.getUser };
    loggedInUser.login = true;
    commit('editUser', loggedInUser);
    // Get boards data from API server...
    return true;
  },
  logout({ commit, getters }) {
    // Dummy logged out process
    const loggedOutUser = { ...getters.getUser };
    loggedOutUser.login = false;
    commit('editUser', loggedOutUser);
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
