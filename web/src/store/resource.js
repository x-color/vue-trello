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
    fetch('/api/resources', {
      headers: {
        'X-XSRF-TOKEN': 'csrf',
        'Content-Type': 'application/json; charset=UTF-8',
      },
      credentials: 'same-origin',
    }).then((response) => {
      if (response.ok) {
        return response.json();
      }
      throw new Error(`Request failed: ${response.status}`);
    }).then((resources) => {
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
