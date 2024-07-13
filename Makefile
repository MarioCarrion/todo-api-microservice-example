GO_VERSION=1.22.5

install:
	go install golang.org/dl/go${GO_VERSION}@latest
	go${GO_VERSION} download
	mkdir -p bin
	ln -sf `go env GOPATH`/bin/go${GO_VERSION} bin/go

lint: tools generate golangci vet dirty

tools:
	go install -C internal/tools \
		github.com/fdaines/spm-go \
		github.com/golangci/golangci-lint/cmd/golangci-lint \
		github.com/jackc/tern/v2 \
		github.com/maxbrunsfeld/counterfeiter/v6 \
		github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen \
		github.com/sqlc-dev/sqlc/cmd/sqlc \
		goa.design/model/cmd/mdl \
		goa.design/model/cmd/stz

generate:
	go generate ./...

golangci:
	golangci-lint run ./...

dirty:
	@status=$$(git status --untracked-files=no --porcelain); \
	if [ ! -z "$${status}" ]; \
	then \
		echo "ERROR: Working directory contains modified files"; \
		git status --untracked-files=no --porcelain; \
		exit 1; \
	fi

vet:
	go vet -tags=redis ./...

test:
	go test -shuffle=on -race -coverprofile=coverage.txt -covermode=atomic $$(go list ./... | grep -v /cmd/)

docker:
	docker compose -f compose.yml -f compose.noop.yml build
	docker compose -f compose.yml -f compose.redis.yml build
	docker compose -f compose.yml -f compose.rabbitmq.yml build
	docker compose -f compose.yml -f compose.kafka.yml build
