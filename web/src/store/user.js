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
        boards: [],
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
  storeBoards(st, user) {
    st.users = st.users.map((u) => {
      if (u.name === user.name) {
        u.boards = user.boards;
      }
      return u;
    });
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
            boards: [],
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
    const i = st.users.findIndex(u => u.name === username && u.password === password);
    if (i !== -1) {
      user.login = true;
      user.boards = st.users[i].boards;
      commit('editUser', user);
      dispatch('loadResources');
    }
    callback(user.login);
  },
  logout({ commit, state: st }) {
    const u = {
      name: st.user.name,
      boards: [...st.user.boards],
    };
    commit('storeBoards', u);

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
