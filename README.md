# OpenTelemetry 入門指南︰建立全面可觀測性架構

---

該儲存庫是 IThome 鐵人賽出版作品 《[OpenTelemetry 入門指南︰建立全面可觀測性架構](https://www.drmaster.com.tw/bookinfo.asp?BookID=MP22334)》的程式範例。

[博碩連結](https://www.drmaster.com.tw/bookinfo.asp?BookID=MP22334)

[天瓏書局連結](https://www.tenlong.com.tw/products/9786263338739?list_name=p-r-zh_tw)

作者︰雷N

[iT邦鐵人檔案](https://ithelp.ithome.com.tw/users/20104930)

[Linkedin](https://www.linkedin.com/in/%E5%81%A5%E8%AA%A0-%E5%91%82-0631b4127/)

`印刷書因圖片太大以及對比問題，能前往各章節圖片資料夾中觀看原圖。`

## 目錄
```
├── 第一部份 : 從傳統到轉型：探索現代 IT 架構
│   ├── CH 1 : 現代化系統架構的演化
│   ├── CH 2 : DevOps 簡介
│   └── CH 3 : 什麼是可觀測性 Observability
├── 第二部份 : 可觀測性開源標準 OpenTelemetry
│   ├── CH 4：淺談 OpenTelemetry
│   ├── CH 5：OpenTelemetry 信號
│   ├── CH 6： Trace SDK 與 Metric SDK介紹
│   ├── CH 7：OpenTelemetry Collector
│   └── CH 8：動手玩OpenTelemetry Demo專案
├── 第三部分：Grafana 開源工具應用
│   ├── CH 9：Grafana 基本概念與應用
│   ├── CH 10：Grafana Loki
│   ├── CH 11：Grafana Tempo
│   └── CH 12：改造 OpenTelemetry Demo專案
├── 第四部分：負載測試工具
│   └── CH 13：Grafana k6 系統測試神器
└── 附錄 Collector 內建指標
```

## 錯字修訂

圖3-12 handnle -> handler

圖4-3 Watefall -> Waterfall

13-15 計量器(Gauge) -> 量規(Gauge)

# 附加文章

Ch 5 

- [OTel Errors](/ch5/OTel%20Errors.md)
- [OTel Go SDK + Slog](/ch5/rolldice/)

Ch 6
- [來玩 OTel Go Metric Exemplar](/ch6/exemplar/)

Ch 7
- [OTel Collector - Filter Processor 的用法很多種？](/ch7/filter_processor.md)

Ch 10
- [Loki Ruler](/ch10/loki_ruler/loki_ruler.md)

Ch13 的目錄中，還有作者最近閱讀 k6 Blog 整理出來的筆記供讀者參考參考。

- [k6 Best Practices and Guidelines](/ch13/k6%20Best%20practices%20and%20guidelines.md)

- [k6 Test Lifecycle](/ch13/k6%20Test%20lifecycle.md)

- [Grafana k6 演示瀏覽器測試](https://youtube.com/playlist?list=PLUPlX-9QUIrMc8ZoG68SWaC54gEnWFa24&si=m1tqV2M5viVTetHY)