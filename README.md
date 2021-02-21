# "ToDo API" Microservice Example

## Introduction

Welcome 👋!

This is an educational repository that includes microservice written in Go. It is used as the principal example of my video series: [Building Microservices in Go/Golang](https://www.youtube.com/playlist?list=PL7yAAGMOat_Fn8sAXIk0WyBfK_sT1pohu).

It's a collection of patterns and guidelines I've successfully used to deliver enterprise microservices when using Go.

The whole purpose of this project is to give you an idea about structuring your Go project with 3 principal goals:

1. It is _enterprise_, meant to last for years,
1. It allows a team to collaborate efficiently with little friction, and
1. It is as idiomatic as possible.

Join the fun at https://youtube.com/MarioCarrion

## Domain Driven Design

This project uses a lot of the ideas introduced by Eric Evans in his book [Domain Driven Design](https://www.domainlanguage.com/), I do encourage reading that book but before I think reading [Domain-Driven Design Distilled](https://smile.amazon.com/Domain-Driven-Design-Distilled-Vaughn-Vernon/dp/0134434420/) makes more sense, also there's a free to download [DDD Reference](https://www.domainlanguage.com/ddd/reference/) available as well.

On YouTube I created [a playlist](https://www.youtube.com/playlist?list=PL7yAAGMOat_GJqfTdM9PBdTRSH7jXs6mI) that includes some of my favorite talks and webminars, feel free to explore that as well.

## Project Structure

Talking specificially about microservices **only**, the structure I like to recommended is the following, everything using `<` and `>` depends on the domain being implemented and the boundext context being defined.

- [ ] `build/`: defines the code used for creating infrastructure
  - [ ] `<cloud-providers>/`: define concrete cloud provider 
  - [ ] `<executableN>/`: contains a Dockerfile used for building the binary
- [ ] `cmd/`
  - [ ] `<primary-server>/`: uses primary database.
  - [ ] `<replica-server>/`: uses readonly databases.
  - [ ] `<binaryN>/`
- [ ] `db/`
  - [ ] `migrations/`: contains database migrations
  - [ ] `seeds/`: contains file meant to populate basic database values
- [ ] `internal/`: defines the _core domain_
  - [ ] `<datastoreN>/`: a concrete _repository_ used by the domain, for example `postgresql`
  - [ ] `http/`: defines HTTP Handlers
  - [ ] `service/`: orchestrates use cases and manages transactions.
- [ ] `pkg/` public API meant to imported by other Go package

## Tools

```
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.14.1
go install github.com/kyleconroy/sqlc/cmd/sqlc@v1.6.0
go install github.com/maxbrunsfeld/counterfeiter/v6@v6.3.0
```

## Features

In no particular order.

- [X] Database migrations
- [X] Repositories
- [X] Dependency Injection
- [X] Secure Configuration
- [ ] Infrastructure as code
- [ ] Metrics and Instrumentation
- [ ] Logging
- [ ] Error Handling
- [ ] Caching
- [ ] Authorization
- [ ] Events
- [ ] Testing
- [ ] Whatever else I forgot to include

## More ideas

* [2016: Peter Bourgon's: Repository structure](https://peter.bourgon.org/go-best-practices-2016/#repository-structure)
* [2016: Ben Johnson's: Standard Package Layout](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1)
* [2017: William Kennedy's: Design Philosophy On Packaging](https://www.ardanlabs.com/blog/2017/02/design-philosophy-on-packaging.html)
* [2017: Jaana Dogan's: Style guideline for Go packages](https://rakyll.org/style-packages/)
* [2018: Kat Zien - How Do You Structure Your Go Apps](https://www.youtube.com/watch?v=oL6JBUk6tj0)

## Docker Containers

### PostgreSQL

Used as repository for persisting data.

```
docker run \
  -d \
  -e POSTGRES_HOST_AUTH_METHOD=trust \
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=dbname \
  -p 5432:5432 \
  postgres:12.5-alpine
```

### Vault

Used as repository for retrieving secrets for configuration values.

```
docker run \
  -d \
  --cap-add=IPC_LOCK \
  -e 'VAULT_DEV_ROOT_TOKEN_ID=myroot' \
  -e 'VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8300' \
  -p 8300:8300 \
  vault:1.6.2
```
