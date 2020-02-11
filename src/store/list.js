import generateUuid from './utils';

// interface List {
//   id: string;
//   boardId: string;
//   title: string;
//   items: []string;
// }

function state() {
  return {
    lists: [],
  };
}

const mutations = {
  addList({ lists }, newList) {
    lists.push(newList);
  },
  removeList(_state, id) {
    _state.lists = _state.lists.filter(list => list.id !== id);
  },
  editList(_state, newList) {
    _state.lists = _state.lists.map((list) => {
      if (list.id === newList.id) {
        return newList;
      }
      return list;
    });
  },
};

const actions = {
  addList({ commit, dispatch, getters }, { title, boardId }) {
    const newList = {
      id: generateUuid(),
      boardId,
      title,
      items: [],
    };
    commit('addList', newList);

    const board = getters.getBoardById(boardId);
    board.lists.push(newList.id);
    dispatch('editBoard', board);
  },
  removeList({ commit, getters, dispatch }, { id, boardId }) {
    getters.getItemsByListId(id).forEach((item) => {
      dispatch('removeItem', item);
    });
    const board = getters.getBoardById(boardId);
    board.lists = board.lists.filter(listId => listId !== id);
    commit('editBoard', board);
    commit('removeList', id);
  },
  editList({ commit }, newList) {
    commit('editList', newList);
  },
};

const getters = {
  getListById: ({ lists }) => id => lists.find(list => list.id === id),
  // eslint-disable-next-line max-len
  getListsByBoardId: (_, _getters) => boardId => _getters.getBoardById(boardId).lists.map(listId => _getters.getListById(listId)),
};

export default {
  state,
  mutations,
  actions,
  getters,
};
