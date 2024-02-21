# README

## Requirements

Install the `tern` tool using [`install_tools`](../bin/install_tools), you can [read more](../internal/tools/) about how those are versioned as well.

## Local PostgreSQL

```
docker run \
  -e POSTGRES_HOST_AUTH_METHOD=trust \
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=dbname \
  -p 5432:5432 \
  postgres:16.2-bullseye
```

## Migrations

Run:

```
tern migrate \
    --migrations "db/migrations/" \
    --conn-string "postgres://user:password@localhost:5432/dbname?sslmode=disable"
```

Create:

```
tern new -m db/migrations/ <migration name>
```
