FROM golang:1.23.0-bookworm AS builder

WORKDIR /build/

COPY db/ .

COPY internal/tools/ .

ENV CGO_ENABLED=1 \
    GOOS=linux

RUN go mod download && go mod tidy

RUN go install -a -installsuffix cgo -ldflags "-extldflags -static" \
		github.com/jackc/tern/v2

#-

FROM debian:bookworm-20240812-slim

ENV DEBIAN_FRONTEND=noninteractive

RUN set -x && \
    apt-get update && \
    apt-get install -y \
      ca-certificates && \
      rm -rf /var/lib/apt/lists/*

WORKDIR /db/
ENV PATH=/db/:$PATH

COPY --from=builder /go/bin/tern .
COPY --from=builder /build/migrations/ .

CMD ["tern"]
