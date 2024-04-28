# CH7：OpenTelemetry Collector

## 環境安裝

1. Docker
2. Docker Compose v2.0.0+

### 啟動服務

```
docker compose up -d
```

### 關閉服務

```
docker compose down
```

## 存取遙測資料的介面

- Grafana︰http://localhost:3000
- Prometheus︰http://localhost:9090
- Collector zPages︰http://localhost:55679/debug/servicez
- Collector pprof︰http://localhost:1888/debug/pprof/
- Collector healthcheck︰http://localhost:13133/

## 透過TelemetryGen產生遙測資料

- 產生日誌

```
docker compose run telemetrygen-logs
```

- 產生指標資料

```
docker compose run telemetrygen-metrics
```

- 產生追蹤跨度

```
docker compose run telemetrygen-traces
```
