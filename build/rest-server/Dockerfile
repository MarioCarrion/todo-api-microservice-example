FROM golang:1.21.5-alpine3.17 AS builder

RUN apk --no-cache add \
      alpine-sdk \
      librdkafka-dev \
      pkgconf && \
    rm -rf /var/cache/apk/*

WORKDIR /build/

COPY . .

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags "-extldflags -static" -tags musl \
      github.com/MarioCarrion/todo-api/cmd/rest-server

RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo -ldflags "-s -w" -tags 'postgres' \
      github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2

#-

FROM alpine:3.18 AS certificates

RUN apk --no-cache add ca-certificates

#-

FROM scratch

WORKDIR /api/
ENV PATH=/api/bin/:$PATH

COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /build/rest-server ./bin/rest-server
COPY --from=builder /go/bin/migrate ./bin/migrate
COPY --from=builder /build/env.example .
COPY --from=builder /build/db/ .

EXPOSE 9234

CMD ["rest-server", "-env", "/api/env.example"]
