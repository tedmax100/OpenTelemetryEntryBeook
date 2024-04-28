# Variables
GO=go
TELEMETRY_GEN=github.com/open-telemetry/opentelemetry-collector-contrib/cmd/telemetrygen@latest
OTLP_ENDPOINT=localhost:4317
OTLP_HTTP_ENDPOINT=localhost:4318

.PHONY: all
all: install-telemetry

.PHONY: install-telemetry
install-telemetry:
	$(GO) install $(TELEMETRY_GEN)

.PHONY: send-traces
send-traces:
	@command -v telemetrygen >/dev/null 2>&1 || { echo >&2 "telemetrygen is not installed. Run 'make install-telemetry' to install it."; exit 1; }
	telemetrygen traces --otlp-insecure --otlp-http $(OTLP_HTTP_ENDPOINT)
	telemetrygen traces --otlp-insecure --otlp-endpoint $(OTLP_ENDPOINT)

.PHONY: send-logs
send-logs:
	@command -v telemetrygen >/dev/null 2>&1 || { echo >&2 "telemetrygen is not installed. Run 'make install-telemetry' to install it."; exit 1; }
	telemetrygen logs --otlp-insecure --otlp-http $(OTLP_HTTP_ENDPOINT)
	telemetrygen logs --otlp-insecure --otlp-endpoint $(OTLP_ENDPOINT)

.PHONY: send-metrics
send-metrics:
	@command -v telemetrygen >/dev/null 2>&1 || { echo >&2 "telemetrygen is not installed. Run 'make install-telemetry' to install it."; exit 1; }
	telemetrygen metrics --otlp-insecure --otlp-http $(OTLP_HTTP_ENDPOINT)
	telemetrygen metrics --otlp-insecure --otlp-endpoint $(OTLP_ENDPOINT)

# docker-compose run telemetrygen traces --duration 5s --otlp-insecure --otlp-endpoint otelcol:4317

.PHONY: gen-data
gen-data:
	@echo "Generating data by calling /rolldice and /register endpoints..."
	@for i in $$(seq 1 20); do \
		curl -s -o /dev/null http://localhost:8080/rolldice || true; \
		curl -s -o /dev/null http://localhost:8080/register || true; \
	done
	@echo "Data generation completed successfully."