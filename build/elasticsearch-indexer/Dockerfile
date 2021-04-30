FROM golang:1.16.3-alpine3.13 AS builder

WORKDIR /build/

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" \
    github.com/MarioCarrion/todo-api/cmd/elasticsearch-indexer

#-

FROM alpine:3.13 AS certificates

RUN apk --no-cache add ca-certificates

#-

FROM scratch

WORKDIR /api/
ENV PATH=/api/bin/:$PATH

COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /build/elasticsearch-indexer ./bin/elasticsearch-indexer
COPY --from=builder /build/env.example .

EXPOSE 9234

CMD ["elasticsearch-indexer", "-env", "/api/env.example"]
