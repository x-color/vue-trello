import generateUuid from './utils';

// interface List {
//   id: string;
//   boardId: string;
//   title: string;
//   index: int;
// }

function state() {
  return {
    lists: [],
  };
}

const mutations = {
  addList({ lists }, boardId, title) {
    lists.push({
      id: generateUuid(),
      boardId,
      title,
      index: lists.length,
    });
  },
  removeList({ lists }, { id }) {
    lists.splice(lists.index(list => list.id === id), 1);
  },
  rearrangeList(_state, oldIndex, newIndex) {
    _state.lists.map((list) => {
      if (list.index === oldIndex) {
        list.index = newIndex;
      } else if (list.index >= newIndex) {
        list.index += 1;
      }
      return list;
    });
  },

};

const actions = {
  addList({ commit }, title) {
    commit('addList', title);
  },
  removeList({ commit }, { id }) {
    commit('removeList', id);
  },
  rearrangeList({ commit }, { index }, newIndex) {
    commit('removeList', index, newIndex);
  },
};

const getters = {
  getListById({ lists }, id) {
    return lists.find(list => list.id === id);
  },
  getListsByBoardId({ lists }, boardId) {
    return lists.filter(list => list.boardId === boardId);
  },
};

export default {
  state,
  mutations,
  actions,
  getters,
};
