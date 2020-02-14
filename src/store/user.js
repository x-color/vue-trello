// interface User {
//   id: string;
//   name: string;
//   boards: []string;
// }

function state() {
  return {
    user: {
      id: '0',
      name: 'testuser',
      login: false,
      background: '',
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
  login({ commit, state: st }, { username, password }) {
    // Dummy login request
    if (!(username === 'testuser' && password === 'password')) {
      return false;
    }
    // Dummy logged in user data
    const loggedInUser = { ...st.user };
    loggedInUser.login = true;
    commit('editUser', loggedInUser);
    // Get boards data from API server...
    return true;
  },
  logout({ commit, state: st }) {
    // Dummy logged out process
    const loggedOutUser = { ...st.user };
    loggedOutUser.login = false;
    commit('editUser', loggedOutUser);
  },
  changeBgColor({ commit, state: st }, { color }) {
    commit('editUser', Object.assign({ ...st.user }, { background: color }));
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
