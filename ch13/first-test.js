// 匯入需要的功能模組
import { check } from 'k6';
import http from 'k6/http';

// 撰寫測試腳本
export default function () {
  const res = http.get('http://test.k6.io/');

  // 撰寫檢查條件
  const checkOutput = check(res, {
    // 檢查收到的 HTTP 狀態是否 200
    'is status 200': (r) => r.status === 200,
    // 檢查收到的內容是否包含標題字串
    'verify homepage text': ( r ) => r.body.includes('Collection of simple web-pages suitable for load testing'),
    // 檢查收到的內容大小是否是 11105 bytes
    'body size is 11,105 bytes': (r) => r.body.length == 11105,
  })

  if (!checkOutput) {
    exit(-1)
  }
}
