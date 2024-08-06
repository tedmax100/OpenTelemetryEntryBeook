.PHONY: build
build:
	docker compose build --no-cache --progress=plain

.PHONY: run
run:
	docker compose up -d

.PHONY: stop
stop:
	docker compose down

.PHONY: load_test
load_test:
	wrk -d30s -c50 http://localhost:8079/hello