# OddJob

[![CI](https://github.com/its-the-vibe/OddJob/actions/workflows/ci.yaml/badge.svg)](https://github.com/its-the-vibe/OddJob/actions/workflows/ci.yaml)

A simple, extensible Go task dispatcher that bridges Redis task messages to the [Poppit](https://github.com/its-the-vibe/Poppit) command execution format.

## Features

- Go service (module-based, no vendoring)
- Config file support (`config.json`) with a tracked example (`config.example.json`)
- Sensitive values from `.env` (tracked `.env.example`, real `.env` ignored)
- Redis integration:
  - `LPOP` task queue messages
  - transform task payloads into Poppit command payloads
  - `RPUSH` into configurable Poppit queue
  - subscribe to Poppit output channel and optionally chain downstream tasks
- Extensible task transformer framework with a `hello:world` reference transformer
- Dockerfile with `scratch` runtime image
- Docker Compose service with `read_only: true`
- Makefile + GitHub Actions CI workflow (`.github/workflows/ci.yaml`)

## Message model

Incoming OddJob task (example):

```json
{
  "taskName": "hello:world",
  "inputFile": "example.txt",
  "metadata": {
    "name": "OddJob"
  }
}
```

The transformer emits a Poppit message:

```json
{
  "repo": "its-the-vibe/OddJob",
  "branch": "refs/heads/main",
  "type": "odd:job",
  "dir": "/workspace",
  "commands": ["echo \"hello OddJob from example.txt\""],
  "metadata": {
    "taskName": "hello:world"
  }
}
```

## Configuration

1. Copy examples:

```bash
cp config.example.json config.json
cp .env.example .env
```

2. Update values in `config.json` (non-sensitive) and `.env` (sensitive values, e.g. `REDIS_PASSWORD`).

3. Run locally:

```bash
make run
```

## Development

```bash
make fmt
make lint
make test
make build
```

## Docker

Build and run with compose (uses external Redis and read-only container filesystem):

```bash
docker compose up --build
```
