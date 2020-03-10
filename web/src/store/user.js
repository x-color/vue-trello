// interface User {
//   id: string;
//   name: string;
//   boards: []string;
// }

function state() {
  return {
    user: {
      name: '',
      login: false,
      color: '',
      boards: [],
    },
    users: [
      {
        name: 'testuser',
        password: 'password',
      },
    ],
  };
}

const mutations = {
  editUser(st, newUser) {
    st.user = newUser;
  },
  addUser(st, newUser) {
    st.users.push(newUser);
  },
};

const actions = {
  editUser({ commit }, newUser) {
    commit('editUser', newUser);
  },
  signup({ commit, state: st }, { username, password, callback }) {
    const res = {
      ok: false,
      status: 400,
    };
    if (typeof username === 'string' && typeof password === 'string') {
      if (username !== '' && password !== '') {
        if (st.users.findIndex(user => user.name === username) === -1) {
          commit('addUser', {
            name: username,
            password,
          });
          res.ok = true;
        } else {
          res.status = 409;
        }
      }
    }
    callback(res);
  },
  login({ commit, dispatch, state: st }, { username, password, callback }) {
    const user = {
      name: username,
      login: false,
      color: '',
      boards: [],
    };
    if (st.users.findIndex(u => u.name === username && u.password === password) !== -1) {
      user.login = true;
      commit('editUser', user);
      dispatch('loadResources');
    }
    callback(user.login);
  },
  logout({ commit }) {
    const user = {
      name: '',
      login: false,
      color: '',
      boards: [],
    };
    commit('editUser', user);
  },
  changeColor({ commit, state: st }, { color }) {
    commit('editUser', Object.assign({ ...st.user }, { color }));
  },
  moveBoard({ commit, state: st }, newBoards) {
    const user = {
      name: st.user.name,
      login: st.user.login,
      color: st.user.color,
      boards: newBoards,
    };
    commit('editUser', user);
  },
};

const getters = {
  user({ user }) {
    return user;
  },
};

export default {
  state,
  mutations,
  actions,
  getters,
};
