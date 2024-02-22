GO_VERSION=1.22.0

tools:
	go install -C internal/tools \
		github.com/deepmap/oapi-codegen/cmd/oapi-codegen \
		github.com/fdaines/spm-go \
		github.com/golangci/golangci-lint/cmd/golangci-lint \
		github.com/jackc/tern/v2 \
		github.com/maxbrunsfeld/counterfeiter/v6 \
		github.com/sqlc-dev/sqlc/cmd/sqlc \
		goa.design/model/cmd/mdl \
		goa.design/model/cmd/stz

install:
	go install golang.org/dl/go${GO_VERSION}@latest
	go${GO_VERSION} download
	mkdir -p bin
	ln -sf `go env GOPATH`/bin/go${GO_VERSION} bin/go

lint: tools generate golangci dirty

dirty:
	@status=$$(git status --untracked-files=no --porcelain); \
	if [ ! -z "$${status}" ]; \
	then \
		echo "ERROR: Working directory contains modified files"; \
		git status --untracked-files=no --porcelain; \
		exit 1; \
	fi

generate:
	go generate ./...

golangci:
	golangci-lint run ./...
