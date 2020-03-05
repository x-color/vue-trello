export default function fetchAPI(url, method = 'GET', body = '') {
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
