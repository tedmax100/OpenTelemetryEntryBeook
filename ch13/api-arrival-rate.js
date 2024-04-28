import http from 'k6/http';

export const options = {
  scenarios: {
    closed_model: {
      executor: 'constant-vus',
      vus: 1,
      duration: '1m',
    },
    open_model: {
      executor: 'constant-arrival-rate',
      rate: 1,
      timeUnit: '1s',
      duration: '1m',
      preAllocatedVUs: 10,
    },
  },
};

export default function () {
  http.get('https://httpbin.test.k6.io/delay/6');
}