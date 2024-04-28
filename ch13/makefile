# Makefile
.PHONY: k6-run

k6-run:
	K6_PROMETHEUS_RW_SERVER_URL=http://localhost:9090/api/v1/write \
	K6_PROMETHEUS_RW_TREND_STATS="p(95),p(99),min,max,med,avg,count,sum" \
	k6 run -o experimental-prometheus-rw $(SCRIPT)

swagger-k6:
	docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli generate \
    -i /local/$(SCRIPT) \
    -g k6 \
    -o /local/${OUT}