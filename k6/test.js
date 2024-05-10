import encoding from 'k6/encoding';
import http from 'k6/http';
import { check } from 'k6';


export let options = {
  iterations: 10000,
  vus: 1000
};

const username = 'user1';
const password = 'password123';

export default function () {
  const credentials = `${username}:${password}`;

  const encodedCredentials = encoding.b64encode(credentials);
  const headerOptions = {
    headers: {
      Authorization: `Basic ${encodedCredentials}`,
      'X-Target-URL': 'http://localhost:8888',
    },
    iterations: 100,
  };

  let res = http.get(`http://localhost:8080/`, headerOptions);

  check(res, {
    'status is 200': (r) => r.status === 200,
    'is authenticated': (r) => r.json().authenticated === true,
    'is correct user': (r) => r.json().user === username,
  });
}
