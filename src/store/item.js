import generateUuid from './utils';

// interface Item {
//   id: string;
//   listId: string;
//   title: string;
//   text: string;
//   tags: Tag[];
// }

function state() {
  return {
    items: [],
  };
}

const mutations = {
  addItem({ items }, {
    listId, title, text, tags,
  }) {
    items.push({
      id: generateUuid(),
      listId,
      title,
      text,
      tags,
    });
  },
  removeItem(_state, id) {
    _state.items = _state.items.filter(item => item.id !== id);
  },
  editItem(_state, newItem) {
    _state.items = _state.items.map((item) => {
      if (item.id === newItem.id) {
        return newItem;
      }
      return item;
    });
  },
};

const actions = {
  addItem({ commit }, {
    listId, title, text = '', tags = [],
  }) {
    commit('addItem', {
      listId, title, text, tags,
    });
  },
  removeItem({ commit }, { id }) {
    commit('removeItem', id);
  },
  editItem({ commit }, newItem) {
    commit('editItem', newItem);
  },
};

const getters = {
  getItemById: ({ items }) => id => items.find(item => item.id === id),
  getItemsByListId: ({ items }) => listId => items.filter(item => item.listId === listId),
};

export default {
  state,
  mutations,
  actions,
  getters,
};
