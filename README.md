# opsctl

A small developer CLI to scaffold Go services and run common local “CI-like” checks.

Built with Go + Cobra. Designed to be simple, fast, and easy to extend.

## Features

- `doctor` — verify required tools are installed (Go, Git, etc.)
- `env check` — validate that `.env` contains all keys from `.env.example`
- `ci run` — run `go fmt`, `go test`, `go build` in a standard order
- `init` — scaffold a minimal Go service template

## Installation

### Option A: Install with Go (recommended for Go developers)

```bash
go install github.com/davidperjans/opsctl/cmd/opsctl@v0.1.3
