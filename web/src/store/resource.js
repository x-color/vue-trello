// interface Tag {
//   id: string;
//   title: string;
//   color: string;
// }

// struct Color {
//   color: string;
// }

function state() {
  const colors = [
    'red',
    'blue',
    'green',
    'yellow',
  ];
  return {
    tags: [
      { id: '0', title: 't1', color: colors[0] },
      { id: '1', title: 't2', color: colors[1] },
      { id: '2', title: 't3', color: colors[2] },
      { id: '3', title: 't4', color: colors[3] },
    ],
    colors,
  };
}

const mutations = {
  setResources(st, { tags, colors }) {
    st.tags = tags;
    st.colors = colors;
  },
};

const actions = {
  // eslint-disable-next-line no-unused-vars
  loadResources({ commit }) {
    return null;
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
