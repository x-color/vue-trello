import fetchAPI from './utils';

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
    fetchAPI('/items', 'POST', JSON.stringify({
      list_id: listId,
      title,
      text,
      tags,
    })).then((item) => {
      const newItem = {
        id: item.id,
        listId: item.list_id,
        title: item.title,
        text: item.text,
        tags: item.tags,
      };
      commit('addItem', newItem);

      const list = getters.getListById(item.list_id);
      list.items.push(item.id);
      commit('editList', list);
    }).catch((err) => {
      console.error(err);
    });
  },
  deleteItem({ commit, getters }, { id, listId }) {
    fetchAPI(`/items/${id}`, 'DELETE').then(() => {
      const list = getters.getListById(listId);
      list.items = list.items.filter(itemId => itemId !== id);
      commit('editList', list);
      commit('deleteItem', id);
    }).catch((err) => {
      console.error(err);
    });
  },
  editItem({ commit }, item) {
    fetchAPI(`/items/${item.id}`, 'PATCH', JSON.stringify({
      list_id: item.listId,
      title: item.title,
      text: item.text,
      tags: item.tags,
    })).then(() => {
      commit('editItem', item);
    }).catch((err) => {
      console.error(err);
    });
  },
  setItems({ commit, state: st }, items) {
    // Remove deleted items from Vuex store
    st.items.filter(item => items.findIndex(i => i.id === item.id) === -1).forEach(item => commit('deleteItem', item));
    // Add or update items
    items.forEach((item) => {
      if (st.items.findIndex(i => i.id === item.id) === -1) {
        commit('addItem', {
          id: item.id,
          listId: item.list_id,
          title: item.title,
          text: item.text,
          tags: item.tags,
        });
      } else {
        commit('editItem', {
          id: item.id,
          listId: item.list_id,
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
