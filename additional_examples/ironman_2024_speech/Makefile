.PHONY: build
build:
	docker compose build --no-cache --progress=plain

.PHONY: run
run:
	docker compose up -d
	@echo ""
	@echo "OpenTelemetry Demo is running."
	@echo "Go to http://localhost:8300 for the Grafana UI."

.PHONY: stop
stop:
	docker compose down

.PHONY: load_test
load_test:
	wrk -d30s -c50 http://localhost:8079/hello