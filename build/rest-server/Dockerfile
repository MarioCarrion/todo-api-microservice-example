FROM golang:1.16.2-alpine3.13 AS builder

WORKDIR /build/

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" \
    github.com/MarioCarrion/todo-api/cmd/rest-server

RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo -ldflags "-s -w" -tags 'postgres' \
    github.com/golang-migrate/migrate/v4/cmd/migrate@v4.14.1

#-

FROM alpine:3.13 AS certificates

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
