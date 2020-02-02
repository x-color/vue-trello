import generateUuid from './utils';

// interface Board {
//   id: string;
//   userId: string;
//   title: string;
//   color: string;
// }

function state() {
  return {
    boards: [{
      id: '1',
      userId: 1,
      title: 'sample board',
      color: 'indigo',
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
    });
  },
  removeBoard({ boards }, { id }) {
    boards.splice(boards.index(board => board.id === id), 1);
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
  removeBoard({ commit }, { id }) {
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
