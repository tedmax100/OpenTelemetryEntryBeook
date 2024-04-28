import http from "k6/http";
import { sleep, check } from "k6";
import { Counter } from "k6/metrics";
import tracing from "k6/experimental/tracing";
import { SharedArray } from "k6/data";


export const CounterErrors = new Counter("error_total_count");

export let options = {
  stages: [
    { duration: "1m", target: 50 },
    { duration: "1m", target: 60 },
    { duration: "1m", target: 70 },
    { duration: "1m", target: 80 },
    { duration: "1m", target: 90 },
    { duration: "1m", target: 100 },
    { duration: "2m", target: 0 },
  ],
  thresholds: {
    // 錯誤總次數 < 100
    error_total_count: [
      {
        threshold: "count<100",
        abortOnFail: true,
      },
    ],
    // 95%的請求在300ms内完成，所有請求在3s内完成
    http_req_duration: [
      {
        threshold: "p(95)<300",
      },
      {
        threshold: "max<3000",
        abortOnFail: true,
      },
    ],
  },
};

var BASE_URL = __ENV.BASE_URL || "http://localhost:8080";

const instrumentedHTTP = new tracing.Client({
  propagator: "w3c",
  sampling: 1.0,
});

// 商品ID
const products = [
  "0PUK6V6EV0",
  "1YMWWN1N4O",
  "2ZYFJ3GM2N",
  "66VCHSJNUP",
  "6E92ZMYYFZ",
  "9SIQT8TOJO",
  "L9ECAV7KIM",
  "LS4PSXUNUM",
  "OLJCESPC7Z",
  "HQTGWGPNH4",
];

// 商品類別
const categories = [
  "binoculars",
  "telescopes",
  "accessories",
  "assembly",
  "travel",
  "books",
  null,
];

// 人員資料
const people = new SharedArray("people data", function () {
  const data = JSON.parse(open("./people.json"));
  return data;
});

export default function () {
  // 訪問首頁
  instrumentedHTTP.request("GET", `${BASE_URL}/`);

  if (Math.random() < 0.7) {
    // 70%的機率查看商品
    let productResponse = instrumentedHTTP.request(
      "GET",
      `${BASE_URL}/api/products/${
        products[Math.floor(Math.random() * products.length)]
      }`,
      null
    );

    if (!check(productResponse, { "is status 200": (r) => r.status === 200 })) {
      CounterErrors.add(1);
      printError(productResponse.request);
    }
  }
  if (Math.random() < 0.2) {
    // 20%的機率查看推薦商品
    let recommendationsResponse = instrumentedHTTP.request(
      "GET",
      `${BASE_URL}/api/recommendations`,
      null,
      { productIds: [randomChoice(products)] }
    );
    if (
      !check(recommendationsResponse, {
        "is status 200": (r) => r.status === 200,
      })
    ) {
      CounterErrors.add(1);
      printError(recommendationsResponse.request);
    }
  }

  // 20%的機率查看廣告
  if (Math.random() < 0.2) {
    http.get(`${BASE_URL}/api/data/`, {
      params: { contextKeys: [randomChoice(categories)] },
    });

    let dataResponse = instrumentedHTTP.request(
      "GET",
      `${BASE_URL}/api/data/`,
      null,
      { contextKeys: [randomChoice(categories)] }
    );
    if (!check(dataResponse, { "is status 200": (r) => r.status === 200 })) {
      CounterErrors.add(1);
      printError(dataResponse.request);
    }
  }

  // 20%的機率查看購物車
  if (Math.random() < 0.2) {
    let cartResponse = instrumentedHTTP.request(
      "GET",
      "http://localhost:8080/api/cart",
      null
    );

    if (!check(cartResponse, { "is status 200": (r) => r.status === 200 })) {
      CounterErrors.add(1);
      printError(cartResponse.request);
    }
  }

  // 15%的機率添加到購物車
  if (Math.random() < 0.15) {
    const product = randomChoice(products);
    const user = __VU;
    const cartItem = {
      item: {
        productId: product,
        quantity: randomChoice([1, 2, 3, 4, 5, 10]),
      },
      userId: user.toString(),
    };

    let addCartResponse = instrumentedHTTP.request(
      "POST",
      `${BASE_URL}/api/cart`,
      JSON.stringify(cartItem),
      { headers: { "Content-Type": "application/json" } }
    );
    if (!check(addCartResponse, { "is status 200": (r) => r.status === 200 })) {
      CounterErrors.add(1);
      printError(addCartResponse.request);
    } else {
      // 5%的機率結帳
      if (Math.random() < 0.5) {
        const checkoutPerson = randomChoice(people);
        const mutableCheckoutPerson = {
          email: checkoutPerson.email,
          address: checkoutPerson.address,
          userCurrency: checkoutPerson.userCurrency,
          creditCard: checkoutPerson.creditCard,
          userId: user.toString(),
        };
        http.post(`${BASE_URL}/api/checkout`, JSON.stringify(checkoutPerson), {
          headers: { "Content-Type": "application/json" },
        });

        let checkoutResponse = instrumentedHTTP.request(
          "POST",
          `${BASE_URL}/api/checkout`,
          JSON.stringify(mutableCheckoutPerson),
          { headers: { "Content-Type": "application/json" } }
        );

        if (
          !check(checkoutResponse, { "is status 200": (r) => r.status === 200 })
        ) {
          CounterErrors.add(1);

          printError(checkoutResponse.request);
        }
      }
    }
  }

  sleep(1);
}

// common functions
// 列印錯誤資訊
function printError(request) {
  console.error(
    `trace_id=${getTraceIdFromTraceparent(request)} method=${
      request.method
    } url=${request.url}`
  );
}

// 從Traceparent header取得trace_id
function getTraceIdFromTraceparent(request) {
  if (!request.headers["Traceparent"]) {
    return "";
  }

  let parts = request.headers["Traceparent"][0].split("-");
  return parts.length > 1 ? parts[1] : "";
}

// 隨機選擇陣列中的元素
function randomChoice(arr) {
  return arr[Math.floor(Math.random() * arr.length)];
}
