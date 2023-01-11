FROM golang:1.19.5-alpine3.17 AS builder

RUN apk --no-cache add \
      alpine-sdk \
      librdkafka-dev \
      pkgconf && \
    rm -rf /var/cache/apk/*

WORKDIR /build/

COPY . .

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags "-extldflags -static" -tags musl \
      -o elasticsearch-indexer github.com/MarioCarrion/todo-api/cmd/elasticsearch-indexer-redis

#-

FROM alpine:3.17 AS certificates

RUN apk --no-cache add ca-certificates

#-

FROM scratch

WORKDIR /api/
ENV PATH=/api/bin/:$PATH

COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /build/elasticsearch-indexer ./bin/elasticsearch-indexer
COPY --from=builder /build/env.example .

CMD ["elasticsearch-indexer", "-env", "/api/env.example"]
