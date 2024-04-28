// import necessary module
import http from 'k6/http';

var testUrl= __ENV.URL ||"https://test-api.k6.io/auth/basic/login/";

export default function () {
  // define URL and payload
  const url = testUrl;
  const payload = JSON.stringify({
    username: 'test_case',
    password: '1234',
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  // send a post request and save response as a variable
  const res = http.post(url, payload, params);
}