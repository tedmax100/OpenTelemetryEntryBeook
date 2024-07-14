---
tags:
  - OTelCollector
  - OTel
  - OTTL
finished: true
---
# OTel Collector - Filter Processor 的用法很多種？

## Filter Processor 常見設定方式
心血來潮打開 TG，看到有大大在詢問關於 OpenTelemetry Collector (contrib) 的 [Filter processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/filterprocessor/README.md) 的使用方法。
第一個問題是關於 Metrics，為什麼官方網站上看到的是 `metrics`︰
```yaml
processors:
    # Data sources: metrics
    filter:
        metrics:
            include:
            match_type: regexp
            metric_names:
                - prefix/.*
                - prefix_.*
```
而套件網站寫的是 `metrics.metric`︰
```yaml
processors:
    filter:
        metrics:
            metric:
                - 'name == "my.metric" and resource.attributes["my_label"] == "abc123"'
                - 'type == METRIC_DATA_TYPE_HISTOGRAM'
            datapoint:
                - 'metric.type == METRIC_DATA_TYPE_SUMMARY'
                - 'resource.attributes["service.name"] == "my_service_name"'
```

第二個問題是關於 Span，一樣官方網站看到的是 `trace.span`︰
```yaml
processors:
  filter/ottl:
    error_mode: ignore
    traces:
      span:
        - |
        resource.attributes["service.name"] != "app1" and
        resource.attributes["service.name"] != "app2" and
        resource.attributes["service.name"] != "app3"

```

看[對岸的部落格](https://wdk-docs.github.io/opentelemetry-docs/docs/collector/transforming-telemetry/#_1)寫的則是 `spans`
```yaml
processors:
  filter/allowlist: # 保留滿足條件的 span
    spans:
      include:
        match_type: strict
        services:
          - app1
          - app2
          - app3
```

問題是什麼樣的寫法是最正確、齊全的用法。

結論是，全部都是正確的寫法！

如果只是要使用正則表達式來查詢`metric name`欄位且為模糊的範圍或該欄位有多種值，查詢會如下。
關於`match_type`，在使用正則表達式需要使用`match_type: regexp`，而如果該欄位的值要正確比對只是有多種值得可能則`match_type: strict`。
又如果是使用內建的 [Filter Metric Label 表達式](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/6a5b9e673a433ad2c6ed7ef2526a8920fd18cf0d/internal/filter/filterexpr/matcher.go#L32)的話，則是使用`match_type: expr`。如`HasLabel`、`Label`。


```yaml
filter/regexp:
    metrics:
        include:
        match_type: regexp
        metric_names: [prefix/.*, .*contains.*, .*_suffix, full_name_match]
filter/regexpincludeoptions:
    metrics: 
        include:
        match_type: regexp
        metric_names:
            - prefix/.*
            - prefix_.*
            - .*contains.*
            - .*_suffix
filter/regexpexcludeoptions:
    metrics: 
        exclude:
            match_type: regex
            metric_names:
                - hello/.*
filter/includestrict:
  metrics:
    include:
      match_type: strict
      metric_names:
        - hello_world
        - hello/world
filter/exprinclude:
  metrics:
    include:
      match_type: expr
      expressions:
        - Label("foo") == "bar"
        - HasLabel("baz")

filter/strict:
    spans:
        include:
            match_type: strict
            services:
            - app1
            - app2
            - app3
            attributes:
            - key: should_include
                value: "(true|probably_true)"
```

但若是要對多個欄位或做更為複雜的查詢操作時會使用到 [OTTL](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/ottl/README.md) ，則會如下
在 Yaml 的表達上就會不同了。且不能正則表達式的比對方法與 OTTL 比對方法同時存在於同一個filter中（`cannot use ottl conditions and include/exclude for metrics at the same time`），原因文章後面會提及。

```yaml
filter/ottldrop:
    error_mode: ignore
    metrics:
      metric:
        - type == METRIC_DATA_TYPE_NONE
    traces:
      span:
        - IsMatch(resource.attributes["k8s.pod.name"], "local")

filter/ottlmetricstraces:
metrics:
    metric: # 針對 HISTOGRAM 類型的資料，name 和 resource.attributes["my_label"] 滿足條件的就過濾
        - 'name == "my.metric" and resource.attributes["my_label"] == "abc123"'
        - 'type == METRIC_DATA_TYPE_HISTOGRAM'
    datapoint: # 針對 datapoint 做過濾，但只針對 metric.type 是 SUMMARY 類型的資料
        - 'metric.type == METRIC_DATA_TYPE_SUMMARY'
traces:
      span: # 針對 Span
        - 'attributes["container.name"] == "app_container_1"'
        - 'resource.attributes["host.name"] == "localhost"'
        - 'name == "app_3"'
      spanevent: # 針對 Span Event
        - 'attributes["grpc"] == true'
        - 'IsMatch(name, ".*grpc.*")'
```

多種混用就要命名多個filter，反正流水線執行順序是依照在各類型遙測資料流水線中安排的順序。記得軟體設計的原則，單一職責的概念，每一個 Filter 盡量只處理單一欄位或比對或單一需求所需的最基本處理。組合排序的部份流給 pipeline。
```yaml
filter/metrics_mix_config:
  metrics:
    include:
      match_type: expr
      expressions:
        - Label("foo") == "bar"
        - HasLabel("baz")

filter/ottlmetricname:
  metric:
    - 'attributes["test"] == "pass"'

filter/spans_mix_config:
  spans:
    include:
      match_type: strict
      services:
        - test
        - test2
      attributes:
        - key: should_include
          value: "(true|probably_true)"
  traces:
    span:
      - 'attributes["test"] == "pass"'
```

> 關於 [OTTL 常見語法](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/ottl/LANGUAGE.md)
> - Paths
> - Lists
> - Literals
> - Enums
> - Converters
> - Math expression

> 關於 [OTTL 常見函數](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/pkg/ottl/ottlfuncs#ottl-functions)
> - append
> - delete
> - keep
> - flatten
> - limit
> - merge_maps
> - replace
> - set
> - truncate


## Filter Processor 程式碼解析

### 設定檔解析
先針對`config.go`擷取部份定義。
設定檔中對應的就是 `Config`，有5種成員屬性，`error_mode`、`metrics`、`logs`、`spans` 和 `traces`。
```go
// Config defines configuration for Resource processor.
type Config struct {
	// ErrorMode determines how the processor reacts to errors that occur while processing an OTTL condition.
	// Valid values are `ignore` and `propagate`.
	// `ignore` means the processor ignores errors returned by conditions and continues on to the next condition. This is the recommended mode.
	// `propagate` means the processor returns the error up the pipeline.  This will result in the payload being dropped from the collector.
	// The default value is `propagate`.
	ErrorMode ottl.ErrorMode `mapstructure:"error_mode"`

	Metrics MetricFilters `mapstructure:"metrics"`

	Logs LogFilters `mapstructure:"logs"`

	Spans filterconfig.MatchConfig `mapstructure:"spans"`

	Traces TraceFilters `mapstructure:"traces"`
}
```

## Metrics Matcher
接著是關於 metrics 區塊中的設定元素。
對應元素有 `include`、`exclude`、`regexp`、`metric` 與 `datapoint`。
```go
// MetricFilters filters by Metric properties.
type MetricFilters struct {
	// Include match properties describe metrics that should be included in the Collector Service pipeline,
	// all other metrics should be dropped from further processing.
	// If both Include and Exclude are specified, Include filtering occurs first.
	Include *filterconfig.MetricMatchProperties `mapstructure:"include"`

	// Exclude match properties describe metrics that should be excluded from the Collector Service pipeline,
	// all other metrics should be included.
	// If both Include and Exclude are specified, Include filtering occurs first.
	Exclude *filterconfig.MetricMatchProperties `mapstructure:"exclude"`

	// RegexpConfig specifies options for the regexp match type
	RegexpConfig *regexp.Config `mapstructure:"regexp"`

	// MetricConditions is a list of OTTL conditions for an ottlmetric context.
	// If any condition resolves to true, the metric will be dropped.
	// Supports `and`, `or`, and `()`
	MetricConditions []string `mapstructure:"metric"`

	// DataPointConditions is a list of OTTL conditions for an ottldatapoint context.
	// If any condition resolves to true, the datapoint will be dropped.
	// Supports `and`, `or`, and `()`
	DataPointConditions []string `mapstructure:"datapoint"`
}
```

接著是關於 span 區塊的設定元素。
對應元素有 `include`、`exclude`。
```go
// MatchConfig has two optional MatchProperties one to define what is processed
// by the processor, captured under the 'include' and the second, exclude, to
// define what is excluded from the processor.
type MatchConfig struct {
	// Include specifies the set of input data properties that must be present in order
	// for this processor to apply to it.
	// Note: If `exclude` is specified, the input data is compared against those
	// properties after the `include` properties.
	// This is an optional field. If neither `include` and `exclude` are set, all input data
	// are processed. If `include` is set and `exclude` isn't set, then all
	// input data matching the properties in this structure are processed.
	Include *MatchProperties `mapstructure:"include"`

	// Exclude specifies when this processor will not be applied to the input data
	// which match the specified properties.
	// Note: The `exclude` properties are checked after the `include` properties,
	// if they exist, are checked.
	// If `include` isn't specified, the `exclude` properties are checked against
	// all input data.
	// This is an optional field. If neither `include` and `exclude` are set, all input data
	// is processed. If `exclude` is set and `include` isn't set, then all the
	// input data that does not match the properties in this structure are processed.
	Exclude *MatchProperties `mapstructure:"exclude"`
}
```

最後是關於 trace 區塊的設定元素。
對應元素有 `span`、`spanevent`。
```go
// TraceFilters filters by OTTL conditions
type TraceFilters struct {
	// SpanConditions is a list of OTTL conditions for an ottlspan context.
	// If any condition resolves to true, the span will be dropped.
	// Supports `and`, `or`, and `()`
	SpanConditions []string `mapstructure:"span"`

	// SpanEventConditions is a list of OTTL conditions for an ottlspanevent context.
	// If any condition resolves to true, the span event will be dropped.
	// Supports `and`, `or`, and `()`
	SpanEventConditions []string `mapstructure:"spanevent"`
}
```

其中 metrics 算是比較複雜的。

### Metrics Filter 程式碼解析

#### MetricFilters
```go
type MetricFilters struct {
    Include              *filterconfig.MetricMatchProperties `mapstructure:"include"`
    Exclude              *filterconfig.MetricMatchProperties `mapstructure:"exclude"`
    RegexpConfig         *regexp.Config                      `mapstructure:"regexp"`
    MetricConditions     []string                            `mapstructure:"metric"`
    DataPointConditions  []string                            `mapstructure:"datapoint"`
}
```
- Include:
包含匹配屬性，描述應包含在流水線中的 metric。
如果指定了 Include 和 Exclude，會首先進行 Include 過濾。
- Exclude:
排除匹配屬性，描述應從流水線中排除的 metric。
如果指定了 Include 和 Exclude，會首先進行 Include 過濾。
- RegexpConfig:
指定正則表達式匹配類型的配置選項。
- MetricConditions:
一組 OTTL 條件，針對 ottlmetric 上下文。如果任何條件為真，則該 metric 將被丟棄。
支援 `and`、`or` 和括號 `()`。
- DataPointConditions:
一組 OTTL 條件，針對 ottldatapoint 上下文。如果任何條件為真，則該 datapoint 將被丟棄。
支援 `and`、`or` 和括號 `()`。


#### MetricMatchProperties
```go
type MetricMatchProperties struct {
    MatchType           MetricMatchType `mapstructure:"match_type"`
    RegexpConfig        *regexp.Config  `mapstructure:"regexp"`
    MetricNames         []string        `mapstructure:"metric_names"`
    Expressions         []string        `mapstructure:"expressions"`
    ResourceAttributes  []Attribute     `mapstructure:"resource_attributes"`
}

type MetricMatchType string

// These are the MetricMatchType that users can specify for filtering
// `pmetric.Metric`s.
const (
	MetricRegexp = MetricMatchType(filterset.Regexp)
	MetricStrict = MetricMatchType(filterset.Strict)
	MetricExpr   = "expr"
)
```
- MatchType:
指定所需的匹配類型（例如 strict 或 regexp）。
- RegexpConfig:
指定 MetricRegexp 匹配類型的配置選項。
- MetricNames:
指定用於匹配 metric 名稱的字串模式列表。如果 metric 名稱匹配此列表中的至少一個字串模式，則匹配成功。
- Expressions:
指定用於匹配 metric 的表達式列表。如果 metric 中的任何 datapoint 匹配此列表中的至少一個表達式，則匹配成功。
- ResourceAttributes:
定義用於匹配 metric 的可能資源屬性列表。如果任何資源屬性匹配此列表中的所有表達式，則匹配成功。

上面的設定配置，被轉成 MetricMatchProperties 後呢就會被載入到 filter processor 的流水線中。
```go
// fsp : filterMetricProcessor
fsp.skipMetricExpr, err = filtermetric.NewSkipExpr(cfg.Metrics.Include, cfg.Metrics.Exclude)
	if err != nil {
		return nil, err
	}
```

然後在 NewSkipExpr 建構式中會看到以下程式碼。會看到是先把 include 加到 matchers 再來才是添加 exclude。所以上述才會說先進行 include 過濾才進行 exclude。

```go
// NewSkipExpr creates a BoolExpr that on evaluation returns true if a metric should NOT be processed or kept.
// The logic determining if a metric should be processed is based on include and exclude settings.
// Include properties are checked before exclude settings are checked.
func NewSkipExpr(include *filterconfig.MetricMatchProperties, exclude *filterconfig.MetricMatchProperties) (expr.BoolExpr[ottlmetric.TransformContext], error) {
	if UseOTTLBridge.IsEnabled() {
		return filterottl.NewMetricSkipExprBridge(include, exclude)
	}
	var matchers []expr.BoolExpr[ottlmetric.TransformContext]
	inclExpr, err := newExpr(include)
	if err != nil {
		return nil, err
	}
	if inclExpr != nil {
		matchers = append(matchers, expr.Not(inclExpr))
	}
	exclExpr, err := newExpr(exclude)
	if err != nil {
		return nil, err
	}
	if exclExpr != nil {
		matchers = append(matchers, exclExpr)
	}
	return expr.Or(matchers...), nil
}

// NewMatcher constructs a metric Matcher. If an 'expr' match type is specified,
// returns an expr matcher, otherwise a name matcher.
func newExpr(mp *filterconfig.MetricMatchProperties) (expr.BoolExpr[ottlmetric.TransformContext], error) {
	if mp == nil {
		return nil, nil
	}
    // 這裡就是 match_type: expr
	if mp.MatchType == filterconfig.MetricExpr {
		if len(mp.Expressions) == 0 {
			return nil, nil
		}
		return newExprMatcher(mp.Expressions)
	}
	if len(mp.MetricNames) == 0 {
		return nil, nil
	}
	return newNameMatcher(mp)
}
```

最後 `match_type: regexp` 的會返回 nameMatcher 類型的matcher。
```go
// nameMatcher matches metrics by metric properties against prespecified values for each property.
type nameMatcher struct {
    // 儲存和處理metric名稱的過濾器集合
	nameFilters filterset.FilterSet
}

func newNameMatcher(mp *filterconfig.MetricMatchProperties) (*nameMatcher, error) {
	nameFS, err := filterset.CreateFilterSet(
		mp.MetricNames,
		&filterset.Config{
			MatchType:    filterset.MatchType(mp.MatchType), // 精確匹配或正則表達式匹配
			RegexpConfig: mp.RegexpConfig, // 正則表達式
		},
	)
	if err != nil {
		return nil, err
	}
	return &nameMatcher{nameFilters: nameFS}, nil
}

// Eval matches a metric using the metric properties configured on the nameMatcher.
// A metric only matches if every metric property configured on the nameMatcher is a match.
func (m *nameMatcher) Eval(_ context.Context, tCtx ottlmetric.TransformContext) (bool, error) {
    // 從上下文中獲取 metric 的名稱來進行比對
	return m.nameFilters.Matches(tCtx.GetMetric().Name()), nil
}
```

<a name="boolexpr"></a>
然後上述的都會轉成 `expr.BoolExpr` 的形式，這說明這比對只會回傳（bool, error），且這比對函數稱為 `Eval`。
```go
// BoolExpr is an interface that allows matching a context K against a configuration of a match.
type BoolExpr[K any] interface {
	Eval(ctx context.Context, tCtx K) (bool, error)
}
```

最後在執行處理 metric 的主要流程, 會在 ScopeMetrics 區段判斷是不是要對該筆資料 skip。
```go
// processMetrics filters the given metrics based off the filterMetricProcessor's filters.
func (fmp *filterMetricProcessor) processMetrics(ctx context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
	if fmp.skipResourceExpr == nil && fmp.skipMetricExpr == nil && fmp.skipDataPointExpr == nil {
		return md, nil
	}

    ...
	var errors error
	md.ResourceMetrics().RemoveIf(func(rmetrics pmetric.ResourceMetrics) bool {
		...
		rmetrics.ScopeMetrics().RemoveIf(func(smetrics pmetric.ScopeMetrics) bool {
			scope := smetrics.Scope()
			smetrics.Metrics().RemoveIf(func(metric pmetric.Metric) bool {
				if fmp.skipMetricExpr != nil {
					skip, err := fmp.skipMetricExpr.Eval(ctx, ottlmetric.NewTransformContext(metric, smetrics.Metrics(), scope, resource, smetrics, rmetrics))
					if err != nil {
						errors = multierr.Append(errors, err)
					}
					if skip {
						return true
					}
				}
				...
			})
			return smetrics.Metrics().Len() == 0
		})
		return rmetrics.ScopeMetrics().Len() == 0
	})

	...
	return md, nil
}
```

#### Config
```go
type Config struct {
    CacheEnabled       bool `mapstructure:"cacheenabled"`
    CacheMaxNumEntries int  `mapstructure:"cachemaxnumentries"`
}
```
- CacheEnabled:
確定是否啟用 LRU 快取以加快後續的匹配速度。如果設置為 true，將啟用快取。
- CacheMaxNumEntries:
LRU 快取中存儲匹配結果的最大資料數量。如果 CacheEnabled 為 false，此設置將被忽略。

到這裡我們應該知道 metrics 中的 include 與 exclude 比對的對象就是 metric_name 以及 resource 中的 attributes，此處大同小異就不重複說明。
regexp 中的 cache 設定是能用很少的記憶體空間，暫存過往結果來加速比對。

接著來看看 `metric` 與 `datapoint`。

## metric 與 datapoint
當設定配置被載入到 filter processor，轉成對應類型的物件後，會呼叫`func (cfg *Config) Validate()`檢查設定配置是否正確。
會看到一句錯誤標語`cannot use ottl conditions and include/exclude for metrics at the same time`；我們不能把 include/exclude 與 OTTL 條件式放在同一個 filter 物件中。所以我們往後看，它這段也只會檢查 OTTL 的文法配置。也呼應最早

```go
// Validate checks if the processor configuration is valid
func (cfg *Config) Validate() error {
	...
	if (cfg.Metrics.MetricConditions != nil || cfg.Metrics.DataPointConditions != nil) && (cfg.Metrics.Include != nil || cfg.Metrics.Exclude != nil) {
		return fmt.Errorf("cannot use ottl conditions and include/exclude for metrics at the same time")
	}
    ...

	var errors error

    ...

	if cfg.Metrics.MetricConditions != nil {
		_, err := filterottl.NewBoolExprForMetric(cfg.Metrics.MetricConditions, filterottl.StandardMetricFuncs(), ottl.PropagateError, component.TelemetrySettings{Logger: zap.NewNop()})
		errors = multierr.Append(errors, err)
	}

	if cfg.Metrics.DataPointConditions != nil {
		_, err := filterottl.NewBoolExprForDataPoint(cfg.Metrics.DataPointConditions, filterottl.StandardDataPointFuncs(), ottl.PropagateError, component.TelemetrySettings{Logger: zap.NewNop()})
		errors = multierr.Append(errors, err)
	}

    ...

	return errors
}
```

再次來看到 filterMetricProcessor 的建構式。能發現這裡與正則表達式用的`newSkipResExpr`不同了。分別是用`filterottl.NewBoolExprForMetric`與`filterottl.NewBoolExprForDataPoint`。

```go
if cfg.Metrics.MetricConditions != nil || cfg.Metrics.DataPointConditions != nil {
		if cfg.Metrics.MetricConditions != nil {
			fsp.skipMetricExpr, err = filterottl.NewBoolExprForMetric(cfg.Metrics.MetricConditions, filterottl.StandardMetricFuncs(), cfg.ErrorMode, set.TelemetrySettings)
			if err != nil {
				return nil, err
			}
		}

		if cfg.Metrics.DataPointConditions != nil {
			fsp.skipDataPointExpr, err = filterottl.NewBoolExprForDataPoint(cfg.Metrics.DataPointConditions, filterottl.StandardDataPointFuncs(), cfg.ErrorMode, set.TelemetrySettings)
			if err != nil {
				return nil, err
			}
		}

		return fsp, nil
	}
```

NewBoolExprForMetric 與 NewBoolExprForDataPoint 這裡就不深入解釋了。反正就是回傳能按照表達式順序依序去執行擷取、轉換、或比對等操作的 OTTL 表達式。
回傳的也是 `BoolExpr` [跳到上述的 boolExpr](#boolexpr)。
```go
// NewBoolExprForMetric creates a BoolExpr[ottlmetric.TransformContext] that will return true if any of the given OTTL conditions evaluate to true.
// The passed in functions should use the ottlmetric.TransformContext.
// If a function named `match` is not present in the function map it will be added automatically so that parsing works as expected
func NewBoolExprForMetric(conditions []string, functions map[string]ottl.Factory[ottlmetric.TransformContext], errorMode ottl.ErrorMode, set component.TelemetrySettings) (expr.BoolExpr[ottlmetric.TransformContext], error) {
	parser, err := ottlmetric.NewParser(functions, set)
	if err != nil {
		return nil, err
	}
	statements, err := parser.ParseConditions(conditions)
	if err != nil {
		return nil, err
	}
	c := ottlmetric.NewConditionSequence(statements, set, ottlmetric.WithConditionSequenceErrorMode(errorMode))
	return &c, nil
}

// NewBoolExprForDataPoint creates a BoolExpr[ottldatapoint.TransformContext] that will return true if any of the given OTTL conditions evaluate to true.
// The passed in functions should use the ottldatapoint.TransformContext.
// If a function named `match` is not present in the function map it will be added automatically so that parsing works as expected
func NewBoolExprForDataPoint(conditions []string, functions map[string]ottl.Factory[ottldatapoint.TransformContext], errorMode ottl.ErrorMode, set component.TelemetrySettings) (expr.BoolExpr[ottldatapoint.TransformContext], error) {
	parser, err := ottldatapoint.NewParser(functions, set)
	if err != nil {
		return nil, err
	}
	statements, err := parser.ParseConditions(conditions)
	if err != nil {
		return nil, err
	}
	c := ottldatapoint.NewConditionSequence(statements, set, ottldatapoint.WithConditionSequenceErrorMode(errorMode))
	return &c, nil
}
```

首先我們要知道 metrics data 到底有哪些屬性。 參考自 [metrics.proto](https://github.com/open-telemetry/opentelemetry-proto/blob/main/opentelemetry/proto/metrics/v1/metrics.proto)。
```
// MetricsData
// └─── ResourceMetrics
//   ├── Resource
//   ├── SchemaURL
//   └── ScopeMetrics
//      ├── Scope
//      ├── SchemaURL
//      └── Metric
//         ├── Name
//         ├── Description
//         ├── Unit
//         └── data
//            ├── Gauge
//            ├── Sum
//            ├── Histogram
//            ├── ExponentialHistogram
//            └── Summary
```

這裡看到的 在 Metric 以下的都是 metric 的範疇。這裡每種 type 裡面的具體監測數值或是 bucket data point 才是屬於 Data point。
像這裡都是 [OTTL metric](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/ottl/contexts/ottlmetric/README.md) 能處理的屬性。最常見的就是  type 也就是上面的 Gauge、Sum、Histogram 等，還有 name 這些。
而 [OTTL DataPoint](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/ottl/contexts/ottldatapoint/README.md) 就不會看到 type 與  name 這種，而是各種 Metric point 內有的數值。有關於 Metric Point 能參考書中的 Ch 5-31 檢測點的說明

舉個例子︰
```yaml
filter/ottlmetricname:
  metrics:
    metric:
      - 'name == "pass"'
    
filter/ottldatapoint:
  metrics:
    datapoint:
      - 'attributes["test"] == "pass"'
      - 'metric.name == "http_client_duration_bucket" and resource.attributes["service.name"] == "notification-sender"'
      - 'start_time_unix_nano < 1720881125000000000'
```

閱讀至此，什麼時候該用 regexp 什麼時候該用 OTTL 來比對 metrics 資料應該有了基本的認知了。


## Trace Matcher 與 Spans Matcher 
用來建立 traces 區塊的是 TraceFilters。
```go
// TraceFilters filters by OTTL conditions
type TraceFilters struct {
	// SpanConditions is a list of OTTL conditions for an ottlspan context.
	// If any condition resolves to true, the span will be dropped.
	// Supports `and`, `or`, and `()`
	SpanConditions []string `mapstructure:"span"`

	// SpanEventConditions is a list of OTTL conditions for an ottlspanevent context.
	// If any condition resolves to true, the span event will be dropped.
	// Supports `and`, `or`, and `()`
	SpanEventConditions []string `mapstructure:"spanevent"`
}
```

而 Spans Matcher 則是 `filterconfig.MatchConfig`。從這組結構我們也能看出與上述的 metric [MetricMatchProperties](#metricmatchproperties)雷同。
```go
// MatchConfig has two optional MatchProperties one to define what is processed
// by the processor, captured under the 'include' and the second, exclude, to
// define what is excluded from the processor.
type MatchConfig struct {
	Include *MatchProperties `mapstructure:"include"`

	Exclude *MatchProperties `mapstructure:"exclude"`
}

// MatchProperties specifies the set of properties in a spans/log/metric to match
// against and if the input data should be included or excluded from the
// processor. At least one of services (spans only), names or
// attributes must be specified. It is supported to have all specified, but
// this requires all the properties to match for the inclusion/exclusion to
// occur.
// The following are examples of invalid configurations:
//
//	attributes/bad1:
//	  # This is invalid because include is specified with neither services or
//	  # attributes.
//	  include:
//	  actions: ...
//
//	span/bad2:
//	  exclude:
//	  	# This is invalid because services, span_names and attributes have empty values.
//	    services:
//	    span_names:
//	    attributes:
//	  actions: ...
//
// Please refer to processor/attributesprocessor/testdata/config.yaml and
// processor/spanprocessor/testdata/config.yaml for valid configurations.
type MatchProperties struct {
	Services []string `mapstructure:"services"`

	SpanNames []string `mapstructure:"span_names"`

	LogBodies []string `mapstructure:"log_bodies"`

	LogSeverityTexts []string `mapstructure:"log_severity_texts"`

	LogSeverityNumber *LogSeverityNumberMatchProperties `mapstructure:"log_severity_number"`

	MetricNames []string `mapstructure:"metric_names"`

	Attributes []Attribute `mapstructure:"attributes"`

	Resources []Attribute `mapstructure:"resources"`

	Libraries []InstrumentationLibrary `mapstructure:"libraries"`

	SpanKinds []string `mapstructure:"span_kinds"`
}
```

能參考書中的 Ch5-62 Trace Span 與 Span context 的關係圖。
Spans Matcher 只能比對關係圖中的 Span 上的屬性，幾乎就是針對名稱或 Kind 或是特定的 trace/span id 做比對。而 Trace Matcher 則能比對到下一層的 Span Context 的內容。

[OTTL Span Context](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/ottl/contexts/ottlspan/README.md)

[OTTL Span Event Context](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/ottl/contexts/ottlspanevent/README.md)

舉個例子︰

```yaml

filter/spans:
  spans:
    include:
      match_type: strict
      services:
        - test
        - test2
      attributes:
        - key: should_include
          value: "(true|probably_true)"
    exclude:
      match_type: regexp
      attributes:
        - key: should_exclude
          value: "(probably_false|false)"

filter/ottl:
  error_mode: ignore
  traces:
    span:
      - 'attributes["test"] == "pass"'
    spanevent:
      - 'attributes["test"] == "pass"'

filter/ottlmultiline:
  traces:
    span:
      - 'attributes["test"] == "pass"'
      - 'attributes["test"] == "also pass"'
```

官網有提供一個易懂的 [Span 資料](https://opentelemetry.io/docs/concepts/signals/traces/#spans)，能讓我們清楚的知道 Span context 與 Span event 的內容。
Span Event 是我們開發者用於紀錄 Span 執行過程中有意義的事件（可能就是業務用，可能是 debug 用）。能參考書中的 Ch5-54 Events 的說明。

```json
{
  "name": "hello-greetings",
  "context": {
    "trace_id": "0x5b8aa5a2d2c872e8321cf37308d69df2",
    "span_id": "0x5fb397be34d26b51"
  },
  "parent_id": "0x051581bf3cb55c13",
  "start_time": "2022-04-29T18:52:58.114304Z",
  "end_time": "2022-04-29T22:52:58.114561Z",
  "attributes": {
    "http.route": "some_route2"
  },
  "events": [
    {
      "name": "hey there!",
      "timestamp": "2022-04-29T18:52:58.114561Z",
      "attributes": {
        "event_attributes": 1
      }
    },
    {
      "name": "bye now!",
      "timestamp": "2022-04-29T18:52:58.114585Z",
      "attributes": {
        "event_attributes": 1
      }
    }
  ]
}

```


# 小結
由於 Filter Processor 與 Attribute Processor 會大量運用到 OTTL 以及正則表達式。但 OTTL 目前還沒到正式釋出 Stable 版本的程度。所以蠻多時候要自行解讀程式碼以及測試案例來理解怎使用。搭配書上的結構關係圖會更好理解。