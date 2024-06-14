Loki use OTLP
====

## What is structured metadata
https://grafana.com/docs/loki/latest/get-started/labels/structured-metadata/

Structured metadata 是一種將metadata 附加到日誌中，而不索引它們或將它們包括在日誌行內容本身中的方法。有用的metadata 示例包括kubernetes pod名稱、處理程序ID或任何其他常用於查詢但具有高基數並且在查詢時提取代價高昂的標籤。
Structured metadata還可以用於從日誌行中查詢常用的metadata，而無需在查詢時應用解析器。例如，大型json blob或使用複雜正則表達式的查詢，會帶來高性能成本。有用的元數據示例包括容器ID或用戶ID。

### 何時使用Structured metadata？
以下情況下應使用Structured metadata：

1. 如果您以OpenTelemetry格式攝取數據，使用Grafana Alloy或OpenTelemetry Collector。Structured metadata旨在支持OpenTelemetry數據的原生攝取。
2. 如果您有不應作為標籤使用且日誌行中不存在的高基數metadata。例子包括進程ID或線程ID或Kubernetes pod名稱。

將已存在於日誌行中的信息提取並放入Structured metadata中是一種反模式。




## Loki distributor use OTLP
https://grafana.com/docs/loki/latest/send-data/otel/