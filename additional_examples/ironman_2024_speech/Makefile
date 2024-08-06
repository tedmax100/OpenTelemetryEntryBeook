.PHONY: node-start
node-start:
	cd internal_service && export OTEL_TRACES_EXPORTER="otlp" && \
	export OTEL_PROPAGATORS="tracecontext,baggage" && \
	export OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318"  && \
	export OTEL_NODE_RESOURCE_DETECTORS="env,host,os" && \
	export OTEL_SERVICE_NAME="internal_service" && \
	export NODE_OPTIONS="--require @opentelemetry/auto-instrumentations-node/register" && \
	node index.js

