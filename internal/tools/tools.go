package tools

import (
	_ "github.com/deepmap/oapi-codegen/cmd/oapi-codegen"       // OpenAPI 3 Code Generator
	_ "github.com/fdaines/spm-go"                              // Software Package Metrics
	_ "github.com/golang-migrate/migrate/v4/cmd/migrate"       // Database Migrations
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Database Migrations
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"    // Linter
	_ "github.com/kyleconroy/sqlc/cmd/sqlc"                    // Type-Safe SQL generator
	_ "github.com/lib/pq"                                      // PostgreSQL Database driver
	_ "github.com/maxbrunsfeld/counterfeiter/v6"               // Mock/Spies/Stubs
	_ "goa.design/model/cmd/mdl"                               // Structurizer
	_ "goa.design/model/cmd/stz"                               // Structurizer
)
