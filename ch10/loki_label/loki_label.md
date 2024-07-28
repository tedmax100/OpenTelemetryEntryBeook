
Loki Stream & Label
====

探討 Loki stream 之前，先來了解 Prometheus 的 Series 這概念。在[官方網站](https://prometheus.io/docs/concepts/data_model/)中這樣說到。
> Prometheus fundamentally stores all data as time series: streams of timestamped values belonging to the same metric and the same set of labeled dimensions. Besides stored time series, Prometheus may generate temporary derived time series as the result of queries.

關鍵點在於`the same metric and the same set of labeled dimensions`。所以會被歸類成同一組 time series 資料的都是同樣的指標名稱以及都是相同的屬性集合。這概念在 Loki 的設計上同樣適用。

![](/ch10/loki_label/img/time-series-2000-opt.png)
[參考自 What is a time series in Prometheus?](https://iximiuz.com/en/posts/prometheus-metrics-labels-time-series/)

可以看到 metric1 對應到 2 個 series，這是因為有一個 label 的不同（instance="foosrv-01:443" 與 instance="foosrv-02:443"）。這就很像關聯式資料庫中的 table，其實這裡就是兩個不同的 table。