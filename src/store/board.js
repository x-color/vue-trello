import generateUuid from './utils';

// interface Board {
//   id: string;
//   userId: string;
//   title: string;
//   index: int;
// }

function state() {
  return {
    boards: [],
  };
}

const mutations = {
  addBoard({ boards }, userId, title) {
    boards.push({
      id: generateUuid(),
      userId,
      title,
      index: boards.length,
    });
  },
  removeBoard({ boards }, { id }) {
    boards.splice(boards.index(board => board.id === id), 1);
  },
  rearrangeBoard(_state, oldIndex, newIndex) {
    _state.boards.map((board) => {
      if (board.index === oldIndex) {
        board.index = newIndex;
      } else if (board.index >= newIndex) {
        board.index += 1;
      }
      return board;
    });
  },
};

const actions = {
  addBoard({ commit }, title) {
    commit('addBoard', title);
  },
  removeBoard({ commit }, { id }) {
    commit('removeBoard', id);
  },
  rearrangeBoard({ commit }, { index }, newIndex) {
    commit('removeBoard', index, newIndex);
  },
};

const getters = {
  getBoardById({ boards }, id) {
    return boards.find(board => board.id === id);
  },
  getBoardsByUserId({ boards }, userId) {
    return boards.filter(board => board.userId === userId);
  },
};

export default {
  state,
  mutations,
  actions,
  getters,
};
