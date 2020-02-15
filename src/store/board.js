import generateUuid from './utils';

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
  removeBoard(st, id) {
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
  removeBoard({ commit, dispatch, getters }, { id }) {
    getters.getListsByBoardId(id).forEach((list) => {
      dispatch('removeList', list);
    });
    const user = { ...getters.user };
    user.boards = user.boards.filter(boardId => boardId !== id);
    commit('editUser', user);
    commit('removeBoard', id);
  },
  editBoard({ commit }, board) {
    commit('editBoard', board);
  },
};

const getters = {
  getBoardById: ({ boards }) => id => boards.find(board => board.id === id),
  // eslint-disable-next-line max-len
  getBoardsByUserId: (_, gtrs) => gtrs.user.boards.map(boardId => gtrs.getBoardById(boardId)),
};

export default {
  state,
  mutations,
  actions,
  getters,
};
