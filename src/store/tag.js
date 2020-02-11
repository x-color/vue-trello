import generateUuid from './utils';

// interface Tag {
//   id: string;
//   title: string;
//   color: string;
// }

function state() {
  return {
    tags: [
      { id: '1', title: 'p1', color: 'red' },
      { id: '2', title: 'p2', color: 'orange' },
      { id: '3', title: 'p3', color: 'green' },
      { id: '4', title: 'p4', color: 'blue' },
    ],
  };
}

const mutations = {
  addTag({ tags }, {
    title, color,
  }) {
    tags.push({
      id: generateUuid(),
      title,
      color,
    });
  },
  removeTag(_state, id) {
    _state.tags = _state.tags.filter(tag => tag.id !== id);
  },
  editTag(_state, newTag) {
    _state.tags = _state.tags.map((tag) => {
      if (tag.id === newTag.id) {
        return newTag;
      }
      return tag;
    });
  },
};

const actions = {
  addTag({ commit }, {
    title, color = 'white',
  }) {
    commit('addTag', {
      title, color,
    });
  },
  removeTag({ commit }, { id }) {
    commit('removeTag', id);
  },
  editTag({ commit }, newTag) {
    commit('editTag', newTag);
  },
};

const getters = {
  getTagById: ({ tags }) => id => tags.find(tag => tag.id === id),
};

export default {
  state,
  mutations,
  actions,
  getters,
};
