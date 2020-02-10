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
  addList({ lists }, { boardId, title }) {
    lists.push({
      id: generateUuid(),
      boardId,
      title,
      items: [],
    });
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
  addList({ commit }, { title, boardId }) {
    commit('addList', { title, boardId });
  },
  removeList({ commit, getters, dispatch }, { id }) {
    getters.getItemsByListId(id).forEach((item) => {
      dispatch('removeItem', item);
    });
    commit('removeList', id);
  },
  editList({ commit }, newList) {
    commit('editList', newList);
  },
};

const getters = {
  getListById: ({ lists }) => id => lists.find(list => list.id === id),
  getListsByBoardId: ({ lists }) => boardId => lists.filter(list => list.boardId === boardId),
};

export default {
  state,
  mutations,
  actions,
  getters,
};
