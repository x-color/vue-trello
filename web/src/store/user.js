import fetchAPI from './utils';

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
  signup(_, { username, password, callback }) {
    const body = JSON.stringify({
      name: username,
      password,
    });

    fetch('/signup', {
      method: 'POST',
      headers: {
        'X-XSRF-TOKEN': 'csrf',
        'Content-Type': 'application/json; charset=UTF-8',
      },
      body,
    }).then((response) => {
      callback(response);
    });
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
  moveBoard({ commit, state: st }, newBoards) {
    const boards = [...st.user.boards];

    const user = {
      name: st.user.name,
      login: st.user.login,
      color: st.user.color,
      boards: newBoards,
    };
    commit('editUser', user);

    boards.some((boardId, i) => {
      if (boardId !== newBoards[i]) {
        let moved;
        let before = '';

        const index = newBoards.indexOf(boardId);
        if (index === i + 1) {
          // Move before 'i'-th board
          // e.g. A B C D E => A B E C D (Moved 'E')
          moved = newBoards[i];
          if (i > 0) {
            before = newBoards[i - 1];
          }
        } else {
          // Move 'i'-th board
          // e.g. A B C D E => A B D E C (Moved 'C')
          moved = boardId;
          before = newBoards[index - 1];
        }

        fetchAPI(`/boards/${moved}/move`, 'PATCH', JSON.stringify({
          before,
        })).catch((err) => {
          console.error(err);
        });

        return true;
      }
      return false;
    });
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
