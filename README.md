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
- Santander statement pipeline:
  - `santander:stmtdate` runs `stmtdate -rename "<file>"`
  - chains to `santander:pdftoppm` using parsed `stmtdate` output
  - `santander:pdftoppm` runs `pdftoppm -png -r 300 "<file>" -f 2 "<file-without-ext>"` in the input file's base directory
  - chains to `santander:stmtpng2tsv`, which runs `${stmtpng2tsv} -input "<file>.png" -output "<file>.tsv"` in the PNG file's directory
  - chains to `santander:stmt2redis`, which runs `${stmt2redis} -f "<file>.tsv" -t santander` in `${basedir}/${orgname}/stmt2redis`
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

The Santander pipeline can be triggered with:

```json
{
  "taskName": "santander:stmtdate",
  "inputFile": "/workspace/incoming/statement.pdf"
}
```

Successful Santander runs chain automatically through `santander:pdftoppm`, `santander:stmtpng2tsv`, and `santander:stmt2redis`.

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

### Aliases

The optional `aliases` section in `config.json` lets you define reusable keys that can be referenced elsewhere in the config or in incoming task messages using the `${key}` syntax:

```json
{
  "aliases": {
    "basedir": "/path/to/directory",
    "orgname": "its-the-vibe",
    "stmtdate": "/path/to/stmtdate",
    "stmtpng2tsv": "/path/to/stmtpng2tsv",
    "stmt2redis": "/path/to/stmt2redis"
  },
  "poppit": {
    "dir": "${basedir}"
  }
}
```

Alias placeholders (`${key}`) are resolved in:

- All `poppit` config fields (`repo`, `branch`, `dir`, `type`, `outputChannel`)
- Transformer-generated Poppit command strings and task-specific working directories
- Incoming task message fields (`inputFile` and `metadata` values)

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
