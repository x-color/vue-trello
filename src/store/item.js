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
  removeItem(st, id) {
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
      listId,
      title,
      text,
      tags,
    };
    commit('addItem', newItem);

    const list = getters.getListById(listId);
    list.items.push(newItem.id);
    commit('editList', list);
  },
  removeItem({ commit, getters }, { id, listId }) {
    const list = getters.getListById(listId);
    list.items = list.items.filter(itemId => itemId !== id);
    commit('editList', list);
    commit('removeItem', id);
  },
  removeItemInDeletedList({ commit }, { id }) {
    commit('removeItem', id);
  },
  editItem({ commit }, newItem) {
    commit('editItem', newItem);
  },
  moveItemAcrossLists({ commit, getters }, { item, toList }) {
    const fromList = { ...getters.getListById(item.listId) };
    fromList.items = fromList.items.filter(itemId => itemId !== item.id);
    item.listId = toList.id;
    commit('editList', fromList);
    commit('editList', toList);
    commit('editItem', item);
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
