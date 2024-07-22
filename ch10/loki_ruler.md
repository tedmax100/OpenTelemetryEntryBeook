Loki Ruler
====

![Loki Components](/各章節圖片/圖10-3.png)
書本中 Ch10 圖10-3 提供了 Loki 每個組件的資料依賴關係。
其中 Ruler ，是 Loki v2.32.0 版本後正式提供的組件。能讓使用者從日誌中通過 LogQL 轉變成指標，可以用來觸發警報，或者評估後用來紀錄的指標。

用來觸發警報的稱為 Alerting Rules，而評估後用來紀錄的指標稱為 Recording Rules。

Grafana 雖然能設定直接對 Loki 設置 Alerting，每隔幾秒鐘就會去執行一次 LogQL 的聚合計算，而且儀表板上也可能很多圖表都是從 Loki 進行聚合近算顯示的結果。聚合計算又通常會搭配 `json`、`count_over_time` 、`unwrap` 或者 `sum` 等對 N 筆日誌進行資料轉型或運算操作。其實都這會耗費 Querier 很多資源。在這些場景下我們要的就只是簡單的指標數值，指標數值的儲存與運算恰好是 Prometheus 擅長的事情，如果我們能將這些資料的儲存與運算的職責交付給它，就能讓 Loki 的資源專注在 Log 的寫入與讀取上。

而且看圖片也能發現 Ruler 並不是對 Querier 進行訪問，而是直接對 Ingester 與持久化儲存進行存取。所以 Ruler 所需的資源就能與 Querier 隔離開來。


## Ruler Recording Rules 的配置
Ruler 簡易配置最主要的有`wal`、`storage`、`rule_path`、`ring`、`enable_api`和`remote_write`。

範例如下︰
```yaml=
ruler:
  wal:
    dir: /tmp/rules/wal
  storage:
    type: local
    lcaol:
      directory: /etc/loki/rules
  rule_path: /tmp/loki/rules
  ring:
    kvstore:
      sotre: inmemory
  enable_api: true
  remote_write:
    enable: true
    client:
      url: http://prometheus:9090/api/v1/write
```

## storage 
用來儲存和存取 Rule 的媒介。內建能使用configdb, azure, gcs, s3, swift,
local, bos, cos。

local 的話就是看自定義的 rules 檔案能從哪裡給 Ruler 取得並加載到記憶體中生成對應的程式物件來執行。

## rule_path
Ruler 有提供 HTTP API 能動態新增或刪減 Rules。

`POST /loki/api/v1/rules/{namespace}`
`DELETE /loki/api/v1/rules/{namespace}/{groupName}`
`DELETE /loki/api/v1/rules/{namespace}`

這些就是會暫時寫入到 rule_path 中。用於臨時新增或實驗用。但最後還是要寫成 Rules 檔案持久化，才能作到版本控制。

## enable_api
則是啟用 Ruler HTTP API。預設是`true`。

## wal
Ruler 計算時生成的暫時資料。當 Ruler 崩潰時，並保證重啟後資料不會遺失，並且可以恢復運行中的狀態。

### dir
指定用於儲存多租戶 WAL 文件的目錄。每個租戶將在此目錄下有自己的子目錄。

### min_age/max_age
設置樣本在WAL中存在的最小/最大時間，只有超過這個時間的樣本才會被截斷。以確保WAL文件不會無限增長並佔用過多儲存空間。

## remote_write
用來將指標寫入到指定的 Prometheus。

### client
Prometheus remote-write endpoint

### add_org_id_header
用於多租戶場景，會在 header 中攜帶 X-Scope-OrgID 的資訊。

## ring
Ruler 配置成 Cluster 場景才有實際作用。有機會在深入講解。

# Rules
Prometheus 也有 Alerting rules的功能，所以[網站上有文件與範例](https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/)。

而 Loki Rules 也是仿照這格式。

### group
```
groups:
  [ - <rule_group> ]
```

### rule_group
```
# The name of the group. Must be unique within a file.
name: <string>

[ interval: <duration> | default = global.evaluation_interval ]

[ limit: <int> | default = 0 ]

[ query_offset: <duration> | default = global.rule_query_offset ]

rules:
  [ - <rule> ... ]
```

## rule
```
# The name of the time series to output to. Must be a valid metric name.
record: <string>

# The PromQL expression to evaluate. 
expr: <string>

# Labels to add or overwrite before storing the result.
labels:
  [ <labelname>: <labelvalue> ]
```

`valid metric name` 這是蠻重要的訊息，不然的話在某些 PromQL 操作時會出現以下警告提示`metric might not be a counter, name does not end in _total/_sum/_bucket`。

能參考 [best practices for naming metrics created by recording rules](https://prometheus.io/docs/practices/rules/#recording-rules)


```yaml
groups:
  - name: NginxRules
    interval: 1m
    rules:
      - record: nginx:requests:rate1m
        expr: |
          sum(
            rate({container="nginx"}[1m])
          )
        labels:
          cluster: "us-central1"
      - record: nginx:error:sum1h
        expr: |
          sum(
            count_over_time({container="nginx"} | json | status >= 400 [1h])
          )
        labels:
          cluster: "us-central1"
```