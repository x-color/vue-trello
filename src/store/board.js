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
    boards: [{
      id: '1',
      userId: 1,
      title: 'sample board',
      text: '',
      color: 'indigo',
      lists: [],
    }],
  };
}

const mutations = {
  addBoard({ boards, userId }, { title, text, color }) {
    boards.push({
      id: generateUuid(),
      userId,
      title,
      text,
      color,
      lists: [],
    });
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
  addBoard({ commit }, { title, text, color }) {
    commit('addBoard', { title, text, color });
  },
  removeBoard({ commit, getters, dispatch }, { id }) {
    getters.getListsByBoardId(id).forEach((list) => {
      dispatch('removeList', list);
    });
    commit('removeBoard', id);
  },
  editBoard({ commit }, board) {
    commit('editBoard', board);
  },
};

const getters = {
  getBoardById: ({ boards }) => id => boards.find(board => board.id === id),
  getBoardsByUserId: ({ boards }) => userId => boards.filter(board => board.userId === userId),
};

export default {
  state,
  mutations,
  actions,
  getters,
};
