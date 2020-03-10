export function fetchAPI(url, method = 'GET', body = '') {
  const options = {
    method,
    headers: {
      'X-XSRF-TOKEN': 'csrf',
      'Content-Type': 'application/json; charset=UTF-8',
    },
    credentials: 'same-origin',
  };
  if (body !== '') {
    options.body = body;
  }

  return fetch(`/api${url}`, options).catch((err) => {
    console.error(err);
    throw new Error(`Request failed: ${err}`);
  }).then((response) => {
    if (response.ok) {
      if (response.status !== 204) {
        return response.json();
      }
      return null;
    }
    throw new Error(`Request failed: ${response.status}`);
  });
}

export function generateUuid() {
  // https://qiita.com/psn/items/d7ac5bdb5b5633bae165
  const chars = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.split('');
  for (let i = 0, len = chars.length; i < len; i += 1) {
    switch (chars[i]) {
      case 'x':
        chars[i] = Math.floor(Math.random() * 16).toString(16);
        break;
      case 'y':
        chars[i] = (Math.floor(Math.random() * 4) + 8).toString(16);
        break;
      default:
        break;
    }
  }
  return chars.join('');
}
