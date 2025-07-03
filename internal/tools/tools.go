package tools

import (
	_ "github.com/fdaines/spm-go"                                // Software Package Metrics
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"      // Linter
	_ "github.com/google/yamlfmt/cmd/yamlfmt"                    // YAML Formatter
	_ "github.com/jackc/tern/v2"                                 // Database Migration
	_ "github.com/maxbrunsfeld/counterfeiter/v6"                 // Mock/Spies/Stubs
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"                        // Type-Safe SQL generator
	_ "golang.org/x/vuln/cmd/govulncheck"                        // Official Go vulnerability checks
)
