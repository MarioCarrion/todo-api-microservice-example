GO_VERSION=1.20.4

tools:
	go install -C internal/tools \
		github.com/deepmap/oapi-codegen/cmd/oapi-codegen \
		github.com/fdaines/spm-go \
		github.com/golangci/golangci-lint/cmd/golangci-lint \
		github.com/jackc/tern/v2 \
		github.com/kyleconroy/sqlc/cmd/sqlc \
		github.com/maxbrunsfeld/counterfeiter/v6 \
		goa.design/model/cmd/mdl \
		goa.design/model/cmd/stz \

install:
	go install golang.org/dl/go${GO_VERSION}@latest
	go${GO_VERSION} download
	mkdir -p bin
	ln -sf `go env GOPATH`/bin/go${GO_VERSION} bin/go
