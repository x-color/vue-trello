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
    },
  };
}

const mutations = {
  editUser(st, newUser) {
    st.user = newUser;
  },
};

const actions = {
  editUser({ commit }, newUser) {
    commit('editUser', newUser);
  },
  login({ commit, dispatch, state: st }, { username, password, callback }) {
    const user = { ...st.user };

    const body = JSON.stringify({
      name: username,
      password,
    });

    fetch('/signin', {
      method: 'POST',
      headers: {
        'X-XSRF-TOKEN': 'csrf',
        'Content-Type': 'application/json; charset=UTF-8',
      },
      body,
    }).then((response) => {
      user.login = response.ok;
      if (user.login) {
        user.name = username;
        commit('editUser', user);
        dispatch('loadResources');
      }
      callback(user.login);
    });
  },
  logout({ commit }) {
    fetch('/signout', {
      headers: {
        'X-XSRF-TOKEN': 'csrf',
        'Content-Type': 'application/json; charset=UTF-8',
      },
    }).then((response) => {
      if (!response.ok) {
        alert('Error: Failed to sign out');
      }
    }).catch(() => {
      alert('Error: Failed to sign out');
    });

    const user = {
      id: '',
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
