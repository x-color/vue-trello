import generateUuid from './utils';

// struct Board {
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
    const newBoard = {
      id: generateUuid(),
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
  // eslint-disable-next-line no-unused-vars
  loadBoards({ commit, getters, state: st }, user) {
    return null;
  },
  // eslint-disable-next-line no-unused-vars
  loadBoard({ commit, dispatch, state: st }, id) {
    return null;
  },
  moveList({ commit }, { board, newLists }) {
    board.lists = newLists;
    commit('editBoard', board);
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
