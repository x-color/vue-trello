import fetchAPI from './utils';

// interface Board {
//   id: string;
//   title: string;
//   text: string;
//   color: string;
//   lists: []string;
// }

function state() {
  return {
    boards: [],
  };
}

const mutations = {
  addBoard({ boards }, newBoard) {
    boards.push(newBoard);
  },
  deleteBoard(st, id) {
    st.boards = st.boards.filter(board => board.id !== id);
  },
  editBoard(st, newBoard) {
    st.boards = st.boards.map((board) => {
      if (board.id === newBoard.id) {
        return newBoard;
      }
      return board;
    });
  },
};

const actions = {
  addBoard({ commit, getters }, {
    title, text, color,
  }) {
    fetchAPI('/boards', 'POST', JSON.stringify({
      title,
      text,
      color,
    })).then((board) => {
      const newBoard = {
        id: board.id,
        title: board.title,
        text: board.text,
        color: board.color,
        lists: [],
      };
      commit('addBoard', newBoard);

      const user = { ...getters.user };
      user.boards.push(newBoard.id);
      commit('editUser', user);
    }).catch((err) => {
      console.error(err);
    });
  },
  deleteBoard({ commit, dispatch, getters }, { id }) {
    fetchAPI(`/boards/${id}`, 'DELETE').then(() => {
      getters.getListsByBoardId(id).forEach((list) => {
        dispatch('deleteListIndeletedBoard', list);
      });
      const user = { ...getters.user };
      user.boards = user.boards.filter(boardId => boardId !== id);
      commit('editUser', user);
      commit('deleteBoard', id);
    }).catch((err) => {
      console.error(err);
    });
  },
  editBoard({ commit }, board) {
    commit('editBoard', board);
  },
  loadBoards({ commit, getters, state: st }, user) {
    fetchAPI('/boards').then(({ boards }) => {
      // Remove deleted boards from Vuex store
      st.boards.forEach(board => commit('deleteBoard', board));
      // Add or update boards
      boards.forEach((board) => {
        if (st.boards.findIndex(b => b.id === board.id) === -1) {
          commit('addBoard', {
            id: board.id,
            title: board.title,
            text: board.text,
            color: board.color,
            lists: [],
          });
        } else {
          commit('editBoard', {
            id: board.id,
            title: board.title,
            text: board.text,
            color: board.color,
            lists: getters.getBoardById(board.id).lists,
          });
        }
      });
      commit('editUser', {
        name: user.name,
        login: user.login,
        color: user.color,
        boards: boards.map(board => board.id),
      });
    }).catch((err) => {
      console.error(err);
    });
  },
  loadBoard({ commit, dispatch, state: st }, id) {
    fetchAPI(`/boards/${id}`).then((board) => {
      if (st.boards.findIndex(b => b.id === board.id) === -1) {
        commit('addBoard', {
          id: board.id,
          title: board.title,
          text: board.text,
          color: board.color,
          lists: board.lists.map(list => list.id),
        });
      } else {
        commit('editBoard', {
          id: board.id,
          title: board.title,
          text: board.text,
          color: board.color,
          lists: board.lists.map(list => list.id),
        });
      }
      dispatch('setLists', board.lists);
    }).catch((err) => {
      console.error(err);
    });
  },
};

const getters = {
  getBoardById: ({ boards }) => id => boards.find(board => board.id === id),
  getSortedBoards: (_, gtrs) => gtrs.user.boards.map(boardId => gtrs.getBoardById(boardId)),
};

export default {
  state,
  mutations,
  actions,
  getters,
};
