import { fetchAPI, generateUuid } from './utils';

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
    // Add list before API request
    const tmpList = {
      id: generateUuid(),
      boardId,
      title,
      items: [],
    };
    commit('addList', tmpList);

    let board = getters.getBoardById(boardId);
    board.lists.push(tmpList.id);
    commit('editBoard', board);

    fetchAPI('/lists', 'POST', JSON.stringify({
      board_id: boardId,
      title,
    })).then((list) => {
      const newList = {
        id: list.id,
        boardId: list.board_id,
        title: list.title,
        items: [],
      };
      commit('addList', newList);

      // Replace temporary list
      board = getters.getBoardById(boardId);
      board.lists = board.lists.map((listId) => {
        if (listId === tmpList.id) {
          return newList.id;
        }
        return listId;
      });
      commit('editBoard', board);
      commit('deleteList', tmpList.id);
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
  deleteListIndeletedBoard({ commit, getters }, { id }) {
    getters.getItemsByListId(id).forEach((item) => {
      commit('deleteItem', item);
    });
    commit('deleteList', id);
  },
  editList({ commit }, list) {
    fetchAPI(`/lists/${list.id}`, 'PATCH', JSON.stringify({
      board_id: list.boardId,
      title: list.title,
    })).then(() => {
      commit('editList', list);
    }).catch((err) => {
      console.error(err);
    });
  },
  setLists({ commit, dispatch, state: st }, lists) {
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
  moveItem({ commit }, { list, newItems }) {
    const items = [...list.items];

    list.items = newItems;
    commit('editList', list);

    if (items.length > newItems.length) {
      // Move item from this list
      return;
    }

    let moved;
    let before = '';

    items.some((itemId, i) => {
      if (itemId !== newItems[i]) {
        const index = newItems.indexOf(itemId);
        if (index === i + 1) {
          // Move before 'i'-th item
          // e.g. A B C D E => A B E C D (Moved 'E')
          moved = newItems[i];
          if (i > 0) {
            before = newItems[i - 1];
          }
        } else {
          // Move 'i'-th item
          // e.g. A B C D E => A B D E C (Moved 'C')
          moved = itemId;
          before = newItems[index - 1];
        }
        return true;
      }
      return false;
    });

    if (newItems.length > items.length && moved === undefined) {
      // Move item to last of this list
      moved = newItems[newItems.length - 1];
      if (newItems.length > 1) {
        before = newItems[newItems.length - 2];
      }
    }

    fetchAPI(`/items/${moved}/move`, 'PATCH', JSON.stringify({
      list_id: list.id,
      before,
    })).catch((err) => {
      console.error(err);
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
