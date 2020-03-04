import { fetchAPI, generateUuid } from './utils';

// interface Board {
//   id: string;
//   userId: string;
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
    userId, title, text, color,
  }) {
    const newBoard = {
      id: generateUuid(),
      userId,
      title,
      text,
      color,
      lists: [],
    };
    commit('addBoard', newBoard);

    const user = { ...getters.user };
    user.boards.push(newBoard.id);
    commit('editUser', user);
  },
  deleteBoard({ commit, dispatch, getters }, { id }) {
    getters.getListsByBoardId(id).forEach((list) => {
      dispatch('deleteListIndeletedBoard', list);
    });
    const user = { ...getters.user };
    user.boards = user.boards.filter(boardId => boardId !== id);
    commit('editUser', user);
    commit('deleteBoard', id);
  },
  editBoard({ commit }, board) {
    commit('editBoard', board);
  },
  loadBoards({ commit, state: st }, user) {
    fetchAPI('/boards').then(({ boards }) => {
      // Remove deleted boards from Vuex store
      st.boards.forEach(board => commit('deleteBoard', board));
      // Add or update boards
      boards.forEach((board) => {
        commit('addBoard', {
          id: board.id,
          userId: board.userId,
          title: board.title,
          text: board.text,
          color: board.color,
          lists: [],
        });
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
          userId: board.user_id,
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
  getBoardsByUserId: (_, gtrs) => gtrs.user.boards.map(boardId => gtrs.getBoardById(boardId)),
};

export default {
  state,
  mutations,
  actions,
  getters,
};
