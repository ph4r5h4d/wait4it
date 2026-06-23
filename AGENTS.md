# AGENTS.md

Guidance for AI agents (and human contributors) working on the **wait4it** codebase.

## Project overview

wait4it is a Go command-line tool that waits until a network port or a backend service
(TCP, MySQL, PostgreSQL, MongoDB, Redis, HTTP, Oracle, InfluxDB, RabbitMQ, Memcached,
Elasticsearch, Aerospike, Kafka, DNS) is ready to accept connections/responses, with an
optional timeout. See [README.md](./README.md) and the published docs at
https://wait4it.dev.

- **Language:** Go (version pinned in [`go.mod`](./go.mod); use `go-version-file: go.mod`).
- **Module path:** `wait4it`
- **Entry point:** [`main.go`](./main.go) (`package main`)
- **License:** see [LICENSE](./LICENSE)

## Repository layout

```
main.go                      # CLI entry point (flags + env parsing)
internal/banner/             # startup banner
pkg/check/                   # check command orchestration + module registry
pkg/model/                   # shared CheckContext struct and Check interface
pkg/<service>/               # one package per supported check (tcp, mysql, postgresql, redis, mongodb, http, oracle, influxdb, rabbitmq, memcached, elasticsearch, aerospike, kafka, dns)
docs/                        # Hugo-based documentation site
Dockerfile / Dockerfile.alpine
Makefile                     # all common operations are defined here
.golangci.yml                # golangci-lint v2 config (errcheck, govet, staticcheck)
.github/workflows/            # CI pipelines (lint, unit-test, integration-test, release, pages)
```

### Per-service package convention

Each `pkg/<service>/` package follows the same pattern:

- `<service>.go` — implementation of the check (implements the `model.Check` interface).
- `<service>_test.go` — **unit tests** (pure, no external services, run by default).
- `<service>_integration_test.go` — **integration tests** guarded by the `//go:build integration` build tag. They use [testcontainers-go](https://github.com/testcontainers/testcontainers-go) and require a running Docker daemon.
- `helpers_test.go` — shared test helpers (where present).

When adding or modifying a check, mirror this structure: add unit tests for the
pure logic and integration tests (under the `integration` build tag) that spin up the
real service via testcontainers.

## Everything must be tested

**Every change must be accompanied by appropriate tests.** Do not ship code without
coverage:

- Add or update **unit tests** for any new logic, helper, parser, or config struct.
  Unit tests must be pure — no network calls, no Docker, no external services — so they
  run reliably in any environment and in CI.
- Add or update **integration tests** (build tag `integration`) for anything that
  touches a real backend. Integration tests use testcontainers, so they require a
  working Docker daemon locally.

Place tests next to the code they cover (`pkg/<service>/`), keeping the
`*_test.go` / `*_integration_test.go` / `helpers_test.go` convention above.

## All gates must pass on every change

Before declaring a task complete, **all** of the following must pass. If any of them
fails, fix the issue before finishing — do not leave the repo in a state where tests
or lint are red.

1. **Unit tests** — `make test-unit`
2. **Integration tests** — `make test-integration` (requires Docker)
3. **Linter** — `make lint` (golangci-lint with `errcheck`, `govet`, `staticcheck`)

To run everything at once, use:

```sh
make test-unit && make lint && make test-integration
```

These three targets mirror exactly what CI runs in
[`.github/workflows/tests.yaml`](./.github/workflows/tests.yaml) (jobs: `lint`,
`unit-test`, `integration-test`). CI will fail on red tests or lint, so the same must
hold locally.

## Always use the Makefile

Run all build, test, and lint operations through the [`Makefile`](./Makefile). Do not
invoke `go test`/`go build`/`golangci-lint` directly with ad-hoc flags — the Makefile
encodes the canonical flags (race detector, `-count=1`, the `integration` build tag,
the installed `golangci-lint` binary, coverage output paths, etc.) and keeps local runs
identical to CI.

Common targets:

| Target | What it does |
| --- | --- |
| `make build` | Build the `wait4it` binary |
| `make run` | Build and run with default flags |
| `make test-unit` | Unit tests only (pure, no Docker) |
| `make test-integration` | Integration tests (build tag `integration`, needs Docker) |
| `make test-all` | Unit + integration tests |
| `make test` | Alias for `test-unit` |
| `make lint` | Run `golangci-lint` |
| `make coverage` | Generate HTML coverage report in `coverage/` |
| `make docker-build` | Build the default (scratch) Docker image |
| `make docker-build-alpine` | Build the Alpine Docker image |
| `make clean` | Remove the binary and coverage artifacts |
| `make help` | List all targets |

If you need a target that doesn't exist yet, add it to the `Makefile` and document it
here, rather than working around it with one-off commands.

## Workflow for a change

1. Reproduce / understand the issue.
2. Write the code change following the per-service package convention.
3. Add or update **unit tests** and **integration tests** as appropriate.
4. Run the gates in order and ensure they are green:
   ```sh
   make test-unit
   make lint
   make test-integration   # Docker must be running
   ```
5. Only once all three pass is the change ready for review/commit.

## Linting notes

- Config lives in [`.golangci.yml`](./golangci.yml) (golangci-lint **v2** format).
- Enabled linters: `errcheck`, `govet`, `staticcheck`.
- Install the linter if missing:
  ```sh
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```
  The `make lint` target will detect a missing binary and print this hint.

## CI

- **Tests** (`.github/workflows/tests.yaml`): runs on PRs and pushes to `main`; executes
  `lint`, `unit-test`, and `integration-test` jobs. Local `make` targets must match
  these jobs.
- **Docker Release** (`.github/workflows/main.yml`): triggered by `v*` tags; builds and
  publishes scratch and alpine images to Docker Hub and GHCR.
- Keep workflow action versions pinned (the repo uses SHA-pinned actions).

## Dependencies & tooling

- Integration tests depend on `github.com/testcontainers/testcontainers-go` and its
  per-service modules (mysql, postgres, redis, mongodb, ...). A running Docker daemon
  is required for `make test-integration`.
- Do not introduce a new service check without also wiring it into
  [`pkg/check/check-module-list.go`](./pkg/check/check-module-list.go) and adding tests.
- Keep `go.mod` / `go.sum` tidy: run `go mod tidy` if module requirements change.