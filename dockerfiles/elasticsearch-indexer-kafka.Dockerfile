FROM golang:1.24.0-bookworm AS builder

WORKDIR /build/

COPY . .

ENV CGO_ENABLED=1 \
    GOOS=linux

RUN go mod download && \
    go build -a -installsuffix cgo -ldflags "-extldflags -static" \
      -o elasticsearch-indexer github.com/MarioCarrion/todo-api-microservice-example/cmd/elasticsearch-indexer-kafka

#-

FROM debian:bookworm-20250203-slim

ENV DEBIAN_FRONTEND=noninteractive

RUN set -x && \
    apt-get update && \
    apt-get install -y --no-install-recommends \
      ca-certificates && \
      rm -rf /var/lib/apt/lists/*

WORKDIR /api/
ENV PATH=/api/bin/:$PATH

COPY --from=builder /build/elasticsearch-indexer ./bin/elasticsearch-indexer
COPY --from=builder /build/env.example .

CMD ["elasticsearch-indexer", "-env", "/api/env.example"]
