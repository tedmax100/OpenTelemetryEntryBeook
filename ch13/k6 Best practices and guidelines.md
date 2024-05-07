---
tags:
  - k6
finished: true
---
# k6 Best practices and guidelines



**k6 的核心價值觀**

1. Treat it as you would any other kind of testing suite.

2. Start simple and then iterate. Basic continuous testing is better than no testing at all.



### **Treat it as you would any other kind of testing suite**

在開發流程中，性能測試通常被視為一個可以延後處理的項目，直到產品出現嚴重的可靠性問題時才被重視。然而，要更加主動地避免這些問題，應該將性能測試視為其他類型測試（如單元測試、整合測試等）一樣的常規工作，從而提前解決潛在的可靠性問題。

**如何將性能測試納入常規測試流程：**

1. **風險優先級**: 在其他類型的測試中，基於失敗風險、關鍵業務流程或問題頻發程度來優先處理測試。同樣地，性能測試時也應該識別哪些功能是最關鍵的，哪些可能會對性能有較大的影響，並對這些部分進行重點測試。

2. **早期測試**: 不要等到產品開發完成或接近完成才開始進行性能測試。性能測試應該從需求分析、設計階段開始就納入考慮，並且在整個開發周期中定期進行，以早期發現問題並處理。

3. **自動化測試**: 將性能測試自動化，使其成為 CI/CD 流程的一部分。這樣可以確保每次代碼更新都會觸發性能測試，及時發現新引入的性能問題。

4. **工具選擇**: 選擇合適的工具如 k6 來執行性能測試。k6 支持腳本化的測試和閾值設定，可以幫助團隊以編程方式定義性能目標並自動化檢測這些目標是否被達到。



**性能測試的特定注意事項：**

- **環境的一致性**: 確保測試環境與營運環境相似，以便測試結果能夠真實反映出應用在營運環境中的表現。

- **多維度測試**: 性能測試不僅僅是測試速度或響應時間，還應該包括系統的可用性、可擴展性和容錯能力等。

- **數據的真實性**: 使用真實世界的數據進行測試，以便測試結果更加準確地反映真實用戶的使用情況。

通過將性能測試視為常規的測試流程的一部分，團隊可以更加主動地發現並解決問題，從而提高產品的質量和用戶的滿意度。這樣的策略不僅提高了測試的覆蓋範圍，還增強了團隊對於性能問題的警覺性和應對能力。

### **Start simple and iterate**

性能測試的逐步演進方法，從簡單開始並逐步迭代擴展。這裡的核心概念和要點可以概括為：



1. **從簡單開始：** 鼓勵從幾個基本測試開始，不需一開始就設置完整的性能測試套件。

2. **經驗與信心的積累：** 隨著團隊對性能測試的熟悉度和信心逐漸增強，測試套件會自然地擴展。

3. **持續測試的重要性：** 與其他類型的測試一樣，性能測試的成功關鍵在於採用持續的測試方法。

4. **進一步的閱讀資源：** 提供了進一步閱讀的建議，即查看[同事Marie Cruz撰寫的文章](https://dzone.com/articles/a-continuous-testing-approach-to-performance)，該文章與k6自動化指南相得益彰，有助於深入了解該主題。

目的在於傳達一種逐步發展和深入理解性能測試的方法，並強調這一過程中持續測試的重要性。





## **It *is* just another testing suite — but often has a broader scope**



性能測試的特點和重要性，強調它與其他測試類型的不同之處以及應對挑戰的方法。這裡的主要觀點和要點可以概括為：

1. **普遍性與特殊性：** 性能測試雖然是一種測試套件，但通常具有更廣泛的範圍。它應該被視為與其他類型的測試一樣，同時需要認識到其獨特的特點。

2. **打破迷思：** 強調性能測試並非困難，困難往往來自於操作的孤立，而非測試實踐本身。

3. **與開發人員緊密合作：** 面對不熟悉系統操作或實現細節的挑戰時，建議與開發人員緊密合作，將[測試左移shift testing left](https://k6.io/modern-load-testing-for-engineering-teams/#work-with-developers-to-shift-testing-left)，採用全團隊的測試方法。

4. **性能測試的簡化實踐：** 對於有測試經驗的開發者來說，性能測試相對簡單，類似於增加了負載維度的自動化功能測試。

5. **迭代測試過程：** 描述了一個測試、開發、再測試的循環過程，直到系統在模擬負載下的表現達到預期的可靠性指標。

6. **不同的測試環境和工作負載：** 強調性能測試通常在開發、QA、預營運、營運等多個環境中進行，並根據不同的流量模式運行不同的負載測試。

7. **負載測試的多樣性：** 受測系統在不同的流量條件下表現不同，進行不同類型的負載測試以驗證預期流量模式下的性能是常見的。

8. **可靠性是團隊努力的結果：** 應用程序的可靠性依賴於底層子系統的可靠性，所有負責這些系統和服務的團隊都應參與測試。

旨在傳達性能測試作為一種獨特測試方式的重要性，以及它需要靈活的測試方法來應對各種挑戰。



## Four best practices for performance testing suites

k6 團隊觀察到常見的效能測試套件的模式。



### 模組化測試配置

性能測試中應用模組化測試配置的最佳實踐，旨在提高測試套件的靈活性和可維護性。以下是主要內容和觀點：

1. **模組化測試的必要性：** 通過模組化，可以在不同的環境和工作負載下重用測試，這提供了更大的靈活性，增強了測試維護，並促進了團隊成員之間的協作。

2. **使用環境變數：** 介紹了如何使用環境變量來指定不同的基本端點。例如，使用 **`BASE_URL`** 來適應不同的測試環境。

3. **儲存環境設定：** 展示了如何使用鍵值對對象儲存每個環境的設定，當環境中有多個設定時，這種方法更為合適。

4. **工作負載配置：** 通過範例顯示如何設置不同的工作負載配置，如平均、stress和smoke測試，以應對不同的流量需求。

5. **靈活性和可擴展性：** 強調了不同環境和工作負載的配置可以根據具體情況進行調整，並可能根據流量不同進行細分。

6. **環境和工作負載的獨立性：** 指出開發或QA環境與營運或預營運環境在資源和可擴展性政策上不同，並且不應在缺乏高可用性設置的環境中運行壓力測試。

7. **配置文件的組織：** 提供了幾種組織配置文件的方法，包括將環境和工作負載設定合併或分開存放。

8. **閾值的模組化：** 展示了如何模組化配置測試的閾值，以便在不同的測試或特定測試群組中重用。

這些描述不僅提供了如何設計和實施性能測試配置的具體指南，也突出了在持續測試過程中保持配置的靈活性和可擴展性的重要性。



#### 使用環境變數

```javascript
const BASE_URL = __ENV.BASE_URL || 'http://localhost:3333';

let res = http.get(`${BASE_URL}/api/pizza`);
```



使用`__ENV`，表示該變數來自環境變數。所以執行上述腳本時只需要 `-e` 搭配變數名稱就能變換變數值。



```javascript
k6 run -e BASE_URL=https://pizza.grafana.fun script.js
```



如果我們該測試腳本需要搭配多個環境時，則能搭配Js的物件來組合。來替換不同環境的BASE_URL。



```javascript
const EnvConfig = {
  dev: {
    BASE_URL: 'http://localhost:3333', 
    MY_FLAG: true
  }
  qa: {
    BASE_URL: 'https://pizza.qa.grafana.fun',
    MY_FLAG: true
  },
  pre: {
    BASE_URL: 'https://pizza.ste.grafana.fun',
    MY_FLAG: false
  },
  prod: {
    BASE_URL: 'https://pizza.grafana.fun',
    MY_FLAG: false
  }
};

const Config = EnvConfig[__ENV.ENVIRONMENT] || EnvConfig['dev'];
const BASE_URL = Config.BASE_URL;
```



執行時只要給定ENVIORMENT的值。

```javascript
k6 run -e ENVIRONMENT=prod script.js
```



又因應第四點，可以配置不同的負載策略。一樣搭配環境變數來更改。



```javascript
const WorkloadConfig = {
  average: [
    { duration: '1m', target: 100 },
    { duration: '4m', target: 100 },
    { duration: '1m', target: 0 },
  ],
  stress: [
    { duration: '1m', target: 700 },
    { duration: '4m', target: 700 },
    { duration: '1m', target: 0 },
  ],
  smoke: [{ duration: '10s', target: 1 }],
};

const stages = WorkloadConfig[__ENV.WORKLOAD] || WorkloadConfig['smoke'];
export const options = {
  stages: stages,
};
```



```javascript
k6 run -e WORKLOAD=stress script.js
```



所以這樣的排列組合下，我們就有 3種負載策略以及4種執行環境，能搭配出12種組合來執行該測試腳本。甚至也能針對不能操作有不同的負載策略，畢竟使用者不會同時使用所有的操作在服務上。



> 也不需要在沒有高可用設計的環境中，進行壓力測試。像是開發整合環境與QA環境通常就不會有高可用/擴縮容等設計，這裡要專注的比較是smoke test這種的功能測試。



> 也不要期待pro-product 環境會跟營運環境有一樣的負載能力。這是不現實的，因為兩者基礎設施等級不同，所以別用一樣的負載與流量對這兩個環境一樣的對待。應該專住在性能回歸測試上，通過測試建立啟性能標準線（baseline），然後頻繁的測試觀察變化找出原因。



考慮上述幾個觀點，所以我們能這樣設計測試腳本根目錄結構。

我們可以把負載策略跟環境設定分開。

```javascript
├── config/
│   ├── workloads.js
│   └── settings.js
├── test1.js
└── test2.js
```



```javascript
// config/workload.js


export const WorkloadConfig = {
  smoke: [...],
  stag: {
    averageLow: [...],
    averageMed: [...],
    averageHigh: [...],
    stress: [...],
  },
  pre: {
    averageLow: [...],
    averageMed: [...],
    averageHigh: [...],
    stress: [...],
    peak: [...],
  },
  prod: {...}
};
```

也能把環境獨立檔案，各環境的負載策略就至於對應環境的檔案中。

```javascript
├── config/
│   ├── dev.js
│   ├── pre.js
│   └── prod.js
├── test1.js
└── test2.js
```







也能把thresholds 閥值給獨立檔案來設定。

```javascript
export const ThresholdsConfig = {
  common: {
      http_req_duration: ['p(99)<1000'],
      http_req_failed: ['rate<0.01']
  },
  // 預營運環境
  pre: {
    instant: {
      http_req_duration: ['p(99)<300'],
    },
  },
  // 營運環境
  prod:{}
};
```





### 實現可重複使用的測試腳本

通常在不同的環境上，我們的測試目標是不同的。例如在QA環境上通常只需要執行smoke test用來進行功能性測試或者確認測試錯誤是否如預期。但同樣的測試腳本執行於預營運環境上，通常是為了驗證SLO目標。而在營運環境上通常我們會使用派程在夜間或系統極少使用者使用時進行測試，以評估性能變化。



![img][圖1]


為了能重複利用測試場景，需要避免將場景邏輯與其他測試概念給緊密的耦合。在設計時能考慮以下觀點︰

- 實現模組化的場景，可以擴展其默認行為

- 讓使用自定義指標成為可選擇項目

- 避免只使用Groups

- 為請求添加Tag以及check 請求結果（以及在適當時添加自定義指標）

- 為了更大的靈活性，使用Scenarios來設置工作負載

- 不要過度思考，從實施意圖的多個測試中重複利用的場景開始。規劃新測試時，一個常見的問題是是否應該擴展現有測試或創建新的測試。



在大多數情況下，我們建議避免多用途測試，並建議為每個場景創建一個新測試，每個環境一個主要目的，如前所述。這樣可以避免混合責任，並有助於追踪同一測試的歷史結果，以識別長期性能變化。



所以我們的目錄結構設計應該很像

```javascript
├── scenarios/
│   ├── e2e/
│   │   ├── checkout.js
│   │   └── read-news.js
│   └── apis/
│       └── account.js
├── tests/
│   ├── smoke-read-news-test.js
│   ├── pre/
│   │   ├── stress-read-news-test.js
│   │   └── avg-read-news-test.js
│   └── prod/
│       └── nightly-read-news-test.js

```



然後就很容易搭配之前的負載策略與環境配置加以組成以下測試腳本。因為每個模組所負責的職責範圍都是獨立的，因此可以相互組合。

```javascript
// tests/pre/avg-read-news-test.js
import { WorkloadConfig, EnvConfig } from './config/workload.js';
import ReadNewsScenario from './scenarios/read-news.js';

const Config = EnvConfig.pre;
const stages = WorkloadConfig.pre.averageMed;

export const options = {
   stages: stages,
}

export default function () {
  ReadNewsScenario(Config.BASE_URL);
}
```



## 用物件導向設計來封裝待測服務的訪問



對待測服務的訪問，我們可以用物件導向設計，將請求的協議跟具體方式都是為一種抽象，所以不管是用HTTP請求還是gRPC請求都可以。將訪問細節封裝到對象中，而對象只提供行為的介面。



下面範例提供一個BaseApiClient的HTTP請求對象，提供HTTP動詞的介面。然後有一個具體的待測服務的API客戶端物件MyApiClient。MyApiClient 提供了待測服務的方法，並藉由BaseApiClient進行擴展。通過這樣的設計，我們就能把數據訪問邏輯風裝在單獨立的類別中，使得測試腳本的程式碼更具有可讀性以及可維護性。



```javascript
class BaseApiClient {
  constructor(baseUrl) {
    this.baseUrl = baseUrl;
  }

  async get(resource) {
    // Logic to perform a GET request
  }

  async post(resource, data) {
    // Logic to perform a POST request
  }

  async put(resource, data) {
    // Logic to perform a PUT request
  }

  async delete(resource) {
    // Logic to perform a DELETE request
  }
}

class MyApiClient extends BaseApiClient {
  constructor(baseUrl) {
    super(baseUrl);
  }

  async getUsers() {
    return this.get('/users');
  }

  async createUser(userData) {
    return this.post('/users', userData);
  }

  async updateUser(userId, userData) {
    return this.put(`/users/${userId}`, userData);
  }

  async deleteUser(userId) {
    return this.delete(`/users/${userId}`);
  }
}

```



## 建立 Error handler wrapper



如果我們送出請求，而沒有處理。則只會看見 `http_req_failed` 會有變化，但不會清楚是哪個測試的請求出錯。雖然測試不會被中斷，但其結果也無法分析利用。

```javascript
http_req_failed................: 36.36% ✓ 40        ✗ 70
```



如果搭配 check 來處理那麼結果就會不同。我們能很明顯的看見是哪個測試對象出現錯誤。

```javascript
check(res, {
   'GET item is 200': (r) => r.status === 200,
});
```

```javascript
✗ GET item is 200
   ↳  63% — ✓ 70 / ✗ 40
```



也能根據同一個測試對象的回應來進行不同的指標分析。

![img][圖2]


傳統上，性能測試主要集中在測試工具本身內部蒐集結果，通常此工具是獨立運行的，這導致了對系統運行方面的可見度不足。



為了能將系統運作方式與其遙測數據，將測試結果與觀察性數據相互關聯。我們能使用[k6/experimental/tracing](https://app.heptabase.com/874af36f-34e9-4f16-9008-ae94afbbfb6f/card/3bbbeded-4b64-41b2-ab5b-b04805519f6d)這模組，它會送出請求時就有一組Trace ID 與 Span ID。方便通過Trace ID來關聯其他遙測資料。



也能通過把錯誤輸出至 Grafana Loki中方便分析。或者創建一個自定義的Counter 指標來紀錄錯誤。而如果我們正在進行高負載測試的場景，可能會面臨失敗數千數萬次，可能會儲存非常大量的數據來分析。但我們可以選擇只儲存端點URL、錯誤訊息、Trace ID或任何相關的資訊即可。



此範例於書本中的 Ch 13 有演示。



以下是一個Error Handler 的演示，如果有錯誤則從中取得Trace ID 以及紀錄上述提及的資訊。

```javascript
class ErrorHandler {
  constructor(logErrorDetails) {
    this.logErrorDetails = logErrorDetails;
  }
  logError(isError, res, tags = {}) {
    if (!isError) return;

    const traceparentHeader = res.request.headers['Traceparent'];
    const errorData = Object.assign(
      {
        url: res.url,
        status: res.status,
        error_code: res.error_code,
        traceparent: traceparentHeader && traceparentHeader.toString(),
      },
      tags
    );
    this.logErrorDetails(errorData);
  }
}
```



然後我們就能選擇是將錯誤輸出至console，還是針對Counter指標+1。

```javascript
const errorHandler = new ErrorHandler((error) => {console.error(error);});

// or

const errors = new CounterMetric('errors');
const errorHandler = new ErrorHandler((error) => {errors.add(1, error);});
```



搭配check 就能進行錯誤處理與輸出。

```javascript

checkStatus = check(res, { 'status is 200': (res) => res.status === 200 });
errorHandler.logError(!checkStatus, res);
```



如果我們選擇將k6的結果通過streaming的形式實時的輸出至儲存服務中，那麼 [Grafana  有提供對應儲存服務的 k6 結果儀表板](https://grafana.com/docs/k6/latest/results-output/grafana-dashboards/)提供展示。





原文參考自 [k6 Best practices and guidelines](https://grafana.com/blog/2024/04/30/organizing-your-grafana-k6-performance-testing-suite-best-practices-to-get-started/?camp=blog&cnt=%F0%9F%92%90+It%27s+gonna+be+May%2C+and&mdm=social&src=li)


[圖1]: ../images/img1.png
[圖2]: ../images/img2.png
