import generateUuid from './utils';

// struct List {
//   id: string;
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
  deleteList(st, id) {
    st.lists = st.lists.filter(list => list.id !== id);
  },
  editList(st, newList) {
    st.lists = st.lists.map((list) => {
      if (list.id === newList.id) {
        return newList;
      }
      return list;
    });
  },
};

const actions = {
  addList({ commit, getters }, { title, boardId }) {
    const newList = {
      id: generateUuid(),
      title,
      items: [],
    };
    commit('addList', newList);

    const board = getters.getBoardById(boardId);
    board.lists.push(newList.id);
    commit('editBoard', board);
  },
  deleteList({ commit, getters }, { id, boardId }) {
    getters.getItemsByListId(id).forEach((item) => {
      commit('deleteItem', item);
    });
    const board = getters.getBoardById(boardId);
    board.lists = board.lists.filter(listId => listId !== id);
    commit('editBoard', board);
    commit('deleteList', id);
  },
  deleteListIndeletedBoard({ commit, getters }, { id }) {
    getters.getItemsByListId(id).forEach((item) => {
      commit('deleteItem', item);
    });
    commit('deleteList', id);
  },
  editList({ commit }, list) {
    commit('editList', list);
  },
  setLists({ commit, dispatch, state: st }, lists) {
    // Add or update lists
    lists.forEach((list) => {
      if (st.lists.findIndex(l => l.id === list.id) === -1) {
        commit('addList', {
          id: list.id,
          title: list.title,
          items: list.items.map(item => item.id),
        });
      } else {
        commit('editList', {
          id: list.id,
          title: list.title,
          items: list.items.map(item => item.id),
        });
      }
      dispatch('setItems', list.items);
    });
  },
  moveItem({ commit }, { list, newItems }) {
    list.items = newItems;
    commit('editList', list);
  },
};

const getters = {
  getListById: ({ lists }) => id => lists.find(list => list.id === id),
  // eslint-disable-next-line max-len
  getListsByBoardId: (_, gtrs) => boardId => gtrs.getBoardById(boardId).lists.map(listId => gtrs.getListById(listId)),
};

export default {
  state,
  mutations,
  actions,
  getters,
};
