import generateUuid from './utils';

// interface Tag {
//   id: string;
//   title: string;
//   color: string;
// }

function state() {
  return {
    tags: [
      { id: generateUuid(), title: 'p1', color: 'red' },
      { id: generateUuid(), title: 'p2', color: 'orange' },
      { id: generateUuid(), title: 'p3', color: 'green' },
      { id: generateUuid(), title: 'p4', color: 'blue' },
    ],
  };
}

const getters = {
  getTagById: ({ tags }) => id => tags.find(tag => tag.id === id),
  tags: ({ tags }) => tags,
};

export default {
  state,
  getters,
};
