import { fetchAPI } from './utils';

// interface Tag {
//   id: string;
//   title: string;
//   color: string;
// }

// interface Color {
//   color: string;
// }

function state() {
  return {
    tags: [],
    colors: [],
  };
}

const mutations = {
  setResources(st, { tags, colors }) {
    st.tags = tags;
    st.colors = colors;
  },
};

const actions = {
  loadResources({ commit }) {
    fetchAPI('/resources').then((resources) => {
      commit('setResources', resources);
    }).catch((err) => {
      console.error(err);
    });
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
