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
  addItem({ items }, newItem) {
    items.push(newItem);
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
  addItem({ commit, dispatch, getters }, {
    listId, title, text = '', tags = [],
  }) {
    const newItem = {
      id: generateUuid(),
      listId,
      title,
      text,
      tags,
    };
    commit('addItem', newItem);

    const list = getters.getListById(listId);
    list.items.push(newItem.id);
    dispatch('editList', list);
  },
  removeItem({ commit, getters }, { id, listId }) {
    const list = getters.getListById(listId);
    list.items = list.items.filter(itemId => itemId !== id);
    commit('editList', list);
    commit('removeItem', id);
  },
  editItem({ commit }, newItem) {
    commit('editItem', newItem);
  },
};

const getters = {
  getItemById: ({ items }) => id => items.find(item => item.id === id),
  // eslint-disable-next-line max-len
  getItemsByListId: (_, _getters) => listId => _getters.getListById(listId).items.map(itemId => _getters.getItemById(itemId)),
};

export default {
  state,
  mutations,
  actions,
  getters,
};
