# migrate

A fork of [golang-migrate/migrate](https://github.com/golang-migrate/migrate) — Database migrations written in Go. Use as CLI or import as library.

[![Go Reference](https://pkg.go.dev/badge/github.com/your-org/migrate.svg)](https://pkg.go.dev/github.com/your-org/migrate)
[![CI](https://github.com/your-org/migrate/actions/workflows/ci.yaml/badge.svg)](https://github.com/your-org/migrate/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/your-org/migrate)](https://goreportcard.com/report/github.com/your-org/migrate)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Stateless** — no external dependency tracking table required by default
- **Multiple database drivers** — PostgreSQL, MySQL, SQLite, and more
- **Multiple source drivers** — filesystem, Go embed, S3, GitHub, and more
- **CLI and library** — use as a standalone tool or import into your Go project
- **Safe** — uses advisory locks to prevent concurrent migrations

## Supported Databases

| Database   | Driver import path                          |
|------------|---------------------------------------------|
| PostgreSQL | `github.com/your-org/migrate/database/postgres` |
| MySQL      | `github.com/your-org/migrate/database/mysql`    |
| SQLite3    | `github.com/your-org/migrate/database/sqlite3`  |
| MongoDB    | `github.com/your-org/migrate/database/mongodb`  |

## Supported Sources

| Source     | Driver import path                          |
|------------|---------------------------------------------|
| File       | `github.com/your-org/migrate/source/file`   |
| Go embed   | `github.com/your-org/migrate/source/iofs`   |
| GitHub     | `github.com/your-org/migrate/source/github` |
| S3         | `github.com/your-org/migrate/source/aws/s3` |

## Installation

### CLI

```bash
# Using Homebrew
brew install your-org/tap/migrate

# Using Go
go install github.com/your-org/migrate/cmd/migrate@latest
```

### Library

```bash
go get github.com/your-org/migrate/v4
```

## Quick Start

### CLI Usage

```bash
# Apply all up migrations
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" up

# Rollback the last migration
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" down 1

# Check current migration version
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" version
```

### Library Usage

```go
package main

import (
    "log"

    "github.com/your-org/migrate/v4"
    _ "github.com/your-org/migrate/v4/database/postgres"
    _ "github.com/your-org/migrate/v4/source/file"
)

func main() {
    m, err := migrate.New(
        "file://./migrations",
        "postgres://localhost:5432/mydb?sslmode=disable",
    )
    if err != nil {
        log.Fatal(err)
    }
    defer m.Close()

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal(err)
    }
}
```

## Migration Files

Migration files follow the naming convention:

```
{version}_{title}.up.{extension}
{version}_{title}.down.{extension}
```

Example:
```
migrations/
  001_create_users.up.sql
  001_create_users.down.sql
  002_add_email_index.up.sql
  002_add_email_index.down.sql
```

## Development

```bash
# Clone the repository
git clone https://github.com/your-org/migrate.git
cd migrate

# Run tests
go test ./...

# Run linter
golangci-lint run

# Build CLI
go build -o migrate ./cmd/migrate
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feat/my-feature`)
3. Commit your changes (`git commit -m 'feat: add my feature'`)
4. Push to the branch (`git push origin feat/my-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

## Acknowledgements

This project is a fork of [golang-migrate/migrate](https://github.com/golang-migrate/migrate). Thanks to all original contributors.
