import generateUuid from './utils';

// interface Item {
//   id: string;
//   listId: string;
//   title: string;
//   text: string;
//   index: int;
// }

function state() {
  return {
    items: [],
  };
}

const mutations = {
  addItem({ items }, listId, title) {
    items.push({
      id: generateUuid(),
      listId,
      title,
      index: items.length,
    });
  },
  removeItem({ items }, { id }) {
    items.splice(items.index(item => item.id === id), 1);
  },
  rearrangeItem(_state, oldIndex, newIndex) {
    _state.items.map((item) => {
      if (item.index === oldIndex) {
        item.index = newIndex;
      } else if (item.index >= newIndex) {
        item.index += 1;
      }
      return item;
    });
  },
};

const actions = {
  addItem({ commit }, title) {
    commit('addItem', title);
  },
  removeItem({ commit }, { id }) {
    commit('removeItem', id);
  },
  rearrangeItem({ commit }, { index }, newIndex) {
    commit('removeItem', index, newIndex);
  },
};

const getters = {
  getItemById({ items }, id) {
    return items.find(item => item.id === id);
  },
  getItemsByListId({ items }, listId) {
    return items.filter(item => item.listId === listId);
  },
};

export default {
  state,
  mutations,
  actions,
  getters,
};
