import generateUuid from './utils';

// interface Item {
//   id: string;
//   listId: string;
//   title: string;
//   text: string;
//   color: string;
// }

function state() {
  return {
    items: [],
  };
}

const mutations = {
  addItem({ items }, {
    listId, title, text = '', color = 'white',
  }) {
    items.push({
      id: generateUuid(),
      listId,
      title,
      text,
      color,
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
    listId, title, text, color,
  }) {
    commit('addItem', {
      listId, title, text, color,
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
