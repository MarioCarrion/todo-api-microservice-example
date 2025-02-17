FROM golang:1.24.0-bookworm AS builder

# Explicitly NOT setting a default value
ARG TAG

WORKDIR /build/

COPY . .

ENV CGO_ENABLED=1 \
    GOOS=linux

RUN go mod download && \
    go build -a -installsuffix cgo -ldflags "-extldflags -static" -tags=$TAG \
		github.com/MarioCarrion/todo-api-microservice-example/cmd/rest-server

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

COPY --from=builder /build/rest-server ./bin/rest-server
COPY --from=builder /build/env.example .

EXPOSE 9234

CMD ["rest-server", "-env", "/api/env.example"]
