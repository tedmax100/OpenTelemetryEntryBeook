import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {

  stages: [
    { duration: '10m', target: 200 }, 
    { duration: '30m', target: 200 }, 
    { duration: '5m', target: 0 }, 
  ],
};

export default () => {
  const urlRes = http.get('https://test-api.k6.io');
  check(urlRes, {
    // 檢查收到的 HTTP 狀態是否 200
    'is status 200': (r) => r.status === 200,
  });
  sleep(1);
};