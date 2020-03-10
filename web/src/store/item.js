import generateUuid from './utils';

// struct Item {
//   id: string;
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
  deleteItem(st, id) {
    st.items = st.items.filter(item => item.id !== id);
  },
  editItem(st, newItem) {
    st.items = st.items.map((item) => {
      if (item.id === newItem.id) {
        return newItem;
      }
      return item;
    });
  },
};

const actions = {
  addItem({ commit, getters }, {
    listId, title, text = '', tags = [],
  }) {
    const newItem = {
      id: generateUuid(),
      title,
      text,
      tags,
    };
    commit('addItem', newItem);

    const list = getters.getListById(listId);
    list.items.push(newItem.id);
    commit('editList', list);
  },
  deleteItem({ commit, getters }, { id, listId }) {
    const list = getters.getListById(listId);
    list.items = list.items.filter(itemId => itemId !== id);
    commit('editList', list);
    commit('deleteItem', id);
  },
  editItem({ commit }, item) {
    commit('editItem', item);
  },
  setItems({ commit, state: st }, items) {
    // Add or update items
    items.forEach((item) => {
      if (st.items.findIndex(i => i.id === item.id) === -1) {
        commit('addItem', {
          id: item.id,
          title: item.title,
          text: item.text,
          tags: item.tags,
        });
      } else {
        commit('editItem', {
          id: item.id,
          title: item.title,
          text: item.text,
          tags: item.tags,
        });
      }
    });
  },
};

const getters = {
  getItemById: ({ items }) => id => items.find(item => item.id === id),
  // eslint-disable-next-line max-len
  getItemsByListId: (_, gtrs) => listId => gtrs.getListById(listId).items.map(itemId => gtrs.getItemById(itemId)),
};

export default {
  state,
  mutations,
  actions,
  getters,
};
