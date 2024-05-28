---
tags:
  - k6
finished: true
---
# k6 Test lifecycle

在 k6 測試的生命週期中，腳本會按以下順序執行：
1. init︰**Required**，每個VU都會執行一次。準備內容階段、設定options、設定全域變數、讀取檔案、載入模組、宣告自定義 function 給VU 程式碼、setup或teardown階段時使用。
2. setup︰測試開始前，執行一次，會在VU之間共享。可以用來呼叫API進行資料初始化。
3. VU 執行的測試程式碼︰**Required**，會是default 函數，或是在scenario函數內。
4. teardown︰每次執行測試腳本時，在測試結束時執行一次

程式碼的結構應該如下︰
```javascript
// 1. init code

export function setup() {
  // 2. setup code
}

export default function (data) {
  // 3. VU code
}

export function teardown(data) {
  // 4. teardown code
}
```


## setup 與 default 函數和 teardown 函數的參數關係
setup 函數返回的數據可以傳遞給 default 和 teardown 函數，方便在不同階段共享數據。
```javascript
import http from 'k6/http';

export function setup() {
  const res = http.get('https://httpbin.test.k6.io/get');
  return { k6res: res.json(), name:"雷N" };
}

export function teardown(data) {
  console.log(JSON.stringify(data));
}

export default function (data) {
  console.log(JSON.stringify(data.k6res), data.name);
}
```