# setup auto-instrumentation with Node.js

## install opentelemetry needed

```shell
pnpm add -S @opentelemetry/api
pnpm add -S @opentelemetry/auto-instrumentations-node
```
## install transporter for winston

```shell
pnpm add -S @opentelemetry/instrumentation-winston
pnpm add -S @opentelemetry/winston-transport
```

## setup arguments Node.js runtime

in dockerfile entrypoint

```yaml
ENTRYPOINT [ "node", "--require", "@opentelemetry/auto-instrumentations-node/register", "main" ]
```
## setup OTEL_EXPORTER_OTLP_ENDPOINT

```shell
export OTEL_EXPORTER_OTLP_ENDPOINT=http://tempo:4318
```

## setup OTEL_NODE_RESOURCE_DETECTORS

```shell
export OTEL_NODE_RESOURCE_DETECTORS=env,host
```
