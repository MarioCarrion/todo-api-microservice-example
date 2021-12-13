# Metrics, Traces and Logging using OpenTelemetry

## Metrics using Prometheus

```
docker run \
  --rm \
  -d \
  -v "${PWD}/prometheus.yml:/etc/prometheus/prometheus.yml" \
  -p 9090:9090 \
  prom/prometheus:v2.31.1
```

Then open http://localhost:9090/

Runtime metrics http://0.0.0.0:9234/metrics

## Tracing using Jaeger

```
docker run \
  --rm \
  -d \
  -p 16686:16686 \
  -p 14268:14268 \
  jaegertracing/all-in-one:1.29.0
```

Then open http://localhost:16686/search
