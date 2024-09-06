# Boilerplate combining Kong CLI argument parsing with zerolog logging

[![pkg.go.dev](https://pkg.go.dev/badge/gitlab.com/tozd/go/cli)](https://pkg.go.dev/gitlab.com/tozd/go/cli)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/tozd/go/cli)](https://goreportcard.com/report/gitlab.com/tozd/go/cli)
[![pipeline status](https://gitlab.com/tozd/go/cli/badges/main/pipeline.svg?ignore_skipped=true)](https://gitlab.com/tozd/go/cli/-/pipelines)
[![coverage report](https://gitlab.com/tozd/go/cli/badges/main/coverage.svg)](https://gitlab.com/tozd/go/cli/-/graphs/main/charts)

A Go package providing boilerplate combining [Kong](https://github.com/alecthomas/kong)
CLI argument parsing with [zerolog](https://gitlab.com/tozd/go/zerolog) logging.

Features:

- Config (from files, CLI arguments, and environment variables) is parsed into a struct
  based on struct tags, powered by [Kong](https://github.com/alecthomas/kong).
- JSON-based and pretty-printed logging, powered by [zerolog](https://gitlab.com/tozd/go/zerolog).
- Support for built-time version variables and CLI flag.
- Handles exit codes: 0 for success, 1 for initialization errors
  (CLI argument parsing or zerolog configuration failures),
  2 for panics, and 3 for program errors.
- All logging goes to stdout and unexpected errors go to stderr.
- On errors, a stack trace and other details from errors are shown,
  powered by [gitlab.com/tozd/go/errors](https://gitlab.com/tozd/go/errors).

## Installation

This is a Go package. You can add it to your project using `go get`:

```sh
go get gitlab.com/tozd/go/cli
```

It requires Go 1.23 or newer.

## Usage

See full package documentation on [pkg.go.dev](https://pkg.go.dev/gitlab.com/tozd/go/cli#section-documentation).

See [examples](./_examples/).

## GitHub mirror

There is also a [read-only GitHub mirror available](https://github.com/tozd/go-cli),
if you need to fork the project there.
