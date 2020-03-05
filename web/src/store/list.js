import fetchAPI from './utils';

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
    fetchAPI('/lists', 'POST', JSON.stringify({
      board_id: boardId,
      title,
    })).then((list) => {
      const newList = {
        id: list.id,
        title: list.title,
        items: [],
      };
      commit('addList', newList);

      const board = getters.getBoardById(boardId);
      board.lists.push(newList.id);
      commit('editBoard', board);
    }).catch((err) => {
      console.error(err);
    });
  },
  deleteList({ commit, getters }, { id, boardId }) {
    fetchAPI(`/lists/${id}`, 'DELETE').then(() => {
      getters.getItemsByListId(id).forEach((item) => {
        commit('deleteItem', item);
      });
      const board = getters.getBoardById(boardId);
      board.lists = board.lists.filter(listId => listId !== id);
      commit('editBoard', board);
      commit('deleteList', id);
    }).catch((err) => {
      console.error(err);
    });
  },
  deleteListIndeletedBoard({ commit, dispatch, getters }, { id }) {
    getters.getItemsByListId(id).forEach((item) => {
      dispatch('deleteItemInDeletedList', item);
    });
    commit('deleteList', id);
  },
  editList({ commit }, newList) {
    commit('editList', newList);
  },
  setLists({ commit, dispatch, state: st }, lists) {
    // Remove deleted lists from Vuex store
    st.lists.filter(list => lists.findIndex(l => l.id === list.id) === -1).forEach(list => commit('deleteList', list));
    // Add or update lists
    lists.forEach((list) => {
      if (st.lists.findIndex(l => l.id === list.id) === -1) {
        commit('addList', {
          id: list.id,
          boardId: list.board_id,
          title: list.title,
          items: list.items.map(item => item.id),
        });
      } else {
        commit('editList', {
          id: list.id,
          boardId: list.board_id,
          title: list.title,
          items: list.items.map(item => item.id),
        });
      }
      dispatch('setItems', list.items);
    });
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
