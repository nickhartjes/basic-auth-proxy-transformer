import encoding from 'k6/encoding';
import http from 'k6/http';
import { check } from 'k6';

const username = 'user';
const password = 'passwd';

export default function () {
  const credentials = `${username}:${password}`;

  const options1 = {
    headers: {
      'X-Target-URL': 'https://httpbin.test.k6.io',
    },
  };


  const url = `https://${credentials}@localhost:8080`;

  let res = http.get(url, options1);

  // Verify response
  check(res, {
    'status is 200': (r) => r.status === 200,
    'is authenticated': (r) => r.json().authenticated === true,
    'is correct user': (r) => r.json().user === username,
  });

  // Alternatively you can create the header yourself to authenticate
  // using HTTP Basic Auth
  const encodedCredentials = encoding.b64encode(credentials);
  const options = {
    headers: {
      Authorization: `Basic ${encodedCredentials}`,
      'X-Target-URL': 'https://httpbin.test.k6.io',
    },
  };

  res = http.get(`https://localhost:8080/basic-auth/${username}/${password}`, options);

  check(res, {
    'status is 200': (r) => r.status === 200,
    'is authenticated': (r) => r.json().authenticated === true,
    'is correct user': (r) => r.json().user === username,
  });
}