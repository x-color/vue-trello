import { fetchAPI, generateUuid } from './utils';

// struct Board {
//   id: string;
//   title: string;
//   text: string;
//   color: string;
//   lists: []string;
// }

function state() {
  return {
    boards: [],
  };
}

const mutations = {
  addBoard({ boards }, newBoard) {
    boards.push(newBoard);
  },
  deleteBoard(st, id) {
    st.boards = st.boards.filter(board => board.id !== id);
  },
  editBoard(st, newBoard) {
    st.boards = st.boards.map((board) => {
      if (board.id === newBoard.id) {
        return newBoard;
      }
      return board;
    });
  },
};

const actions = {
  addBoard({ commit, getters }, {
    title, text, color,
  }) {
    // Add board before API request
    const tmpBoard = {
      id: generateUuid(),
      title,
      text,
      color,
      lists: [],
    };
    commit('addBoard', tmpBoard);

    let user = { ...getters.user };
    user.boards.push(tmpBoard.id);
    commit('editUser', user);

    fetchAPI('/boards', 'POST', JSON.stringify({
      title,
      text,
      color,
    })).then((board) => {
      const newBoard = {
        id: board.id,
        title: board.title,
        text: board.text,
        color: board.color,
        lists: [],
      };
      commit('addBoard', newBoard);

      // Replace temporary board
      user = { ...getters.user };
      user.boards = user.boards.map((boardId) => {
        if (boardId === tmpBoard.id) {
          return newBoard.id;
        }
        return boardId;
      });
      commit('editUser', user);
      commit('deleteBoard', tmpBoard.id);
    }).catch((err) => {
      console.error(err);
    });
  },
  deleteBoard({ commit, dispatch, getters }, { id }) {
    fetchAPI(`/boards/${id}`, 'DELETE').then(() => {
      getters.getListsByBoardId(id).forEach((list) => {
        dispatch('deleteListIndeletedBoard', list);
      });
      const user = { ...getters.user };
      user.boards = user.boards.filter(boardId => boardId !== id);
      commit('editUser', user);
      commit('deleteBoard', id);
    }).catch((err) => {
      console.error(err);
    });
  },
  editBoard({ commit }, board) {
    fetchAPI(`/boards/${board.id}`, 'PATCH', JSON.stringify({
      title: board.title,
      text: board.text,
      color: board.color,
    })).then(() => {
      commit('editBoard', board);
    }).catch((err) => {
      console.error(err);
    });
  },
  loadBoards({ commit, getters, state: st }, user) {
    fetchAPI('/boards').then(({ boards }) => {
      // Add or update boards
      boards.forEach((board) => {
        if (st.boards.findIndex(b => b.id === board.id) === -1) {
          commit('addBoard', {
            id: board.id,
            title: board.title,
            text: board.text,
            color: board.color,
            lists: [],
          });
        } else {
          commit('editBoard', {
            id: board.id,
            title: board.title,
            text: board.text,
            color: board.color,
            lists: getters.getBoardById(board.id).lists,
          });
        }
      });
      commit('editUser', {
        name: user.name,
        login: user.login,
        color: user.color,
        boards: boards.map(board => board.id),
      });
    }).catch((err) => {
      console.error(err);
    });
  },
  loadBoard({ commit, dispatch, state: st }, id) {
    fetchAPI(`/boards/${id}`).then((board) => {
      if (st.boards.findIndex(b => b.id === board.id) === -1) {
        commit('addBoard', {
          id: board.id,
          title: board.title,
          text: board.text,
          color: board.color,
          lists: board.lists.map(list => list.id),
        });
      } else {
        commit('editBoard', {
          id: board.id,
          title: board.title,
          text: board.text,
          color: board.color,
          lists: board.lists.map(list => list.id),
        });
      }
      dispatch('setLists', board.lists);
    }).catch((err) => {
      console.error(err);
    });
  },
  moveList({ commit }, { board, newLists }) {
    const lists = [...board.lists];

    board.lists = newLists;
    commit('editBoard', board);

    lists.some((listId, i) => {
      if (listId !== newLists[i]) {
        let moved;
        let before = '';

        const index = newLists.indexOf(listId);
        if (index === i + 1) {
          // Move before 'i'-th list
          // e.g. A B C D E => A B E C D (Moved 'E')
          moved = newLists[i];
          if (i > 0) {
            before = newLists[i - 1];
          }
        } else {
          // Move 'i'-th list
          // e.g. A B C D E => A B D E C (Moved 'C')
          moved = listId;
          before = newLists[index - 1];
        }

        fetchAPI(`/lists/${moved}/move`, 'PATCH', JSON.stringify({
          board_id: board.id,
          before,
        })).catch((err) => {
          console.error(err);
        });

        return true;
      }
      return false;
    });
  },
};

const getters = {
  getBoardById: ({ boards }) => id => boards.find(board => board.id === id),
  getSortedBoards: (_, gtrs) => gtrs.user.boards.map(boardId => gtrs.getBoardById(boardId)),
};

export default {
  state,
  mutations,
  actions,
  getters,
};
