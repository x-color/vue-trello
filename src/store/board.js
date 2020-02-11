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
  removeBoard(_state, id) {
    _state.boards = _state.boards.filter(board => board.id !== id);
  },
  editBoard(_state, newBoard) {
    _state.boards = _state.boards.map((board) => {
      if (board.id === newBoard.id) {
        return newBoard;
      }
      return board;
    });
  },
};

const actions = {
  addBoard({ commit, dispatch, getters }, {
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

    const user = getters.getUser;
    user.boards.push(newBoard.id);
    dispatch('editUser', user);
  },
  removeBoard({ commit, getters, dispatch }, { id }) {
    getters.getListsByBoardId(id).forEach((list) => {
      dispatch('removeList', list);
    });
    const user = getters.getUser;
    user.boards = user.boards.filter(boardId => boardId !== id);
    dispatch('editUser', user);
    commit('removeBoard', id);
  },
  editBoard({ commit }, board) {
    commit('editBoard', board);
  },
};

const getters = {
  getBoardById: ({ boards }) => id => boards.find(board => board.id === id),
  // eslint-disable-next-line max-len
  getBoardsByUserId: (_, _getters) => _getters.getUser.boards.map(boardId => _getters.getBoardById(boardId)),
};

export default {
  state,
  mutations,
  actions,
  getters,
};
