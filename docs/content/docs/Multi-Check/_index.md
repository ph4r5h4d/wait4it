+++
date = '2026-06-24T00:00:00+01:00'
draft = false
title = 'Multi-Check YAML Configuration'
weight = 15
+++
# Multi-Check YAML Configuration

In addition to single checks via flags, Wait4it supports running multiple checks from a YAML file (via `--config` / `W4IT_CONFIG`) or inline YAML content (via `W4IT_CONFIG_YAML` environment variable).

When a configuration source is provided, all other check-specific flags and their environment variable equivalents are ignored.

## Basic Usage

### Using a YAML file

```bash
wait4it --config checks.yaml
# or
W4IT_CONFIG=checks.yaml wait4it
```

### Using inline YAML

```bash
W4IT_CONFIG_YAML='version: "1"
checks:
  - type: tcp
    host: 127.0.0.1
    port: 8080' wait4it
```

## YAML Structure

```yaml
version: "1"          # Optional, only "1" supported
timeout: 30           # Global default timeout in seconds
fail_fast: true       # Stop on first required failure (default: true)

checks:               # List of checks (required, non-empty)
  - name: "primary-db"   # Optional human-readable name for logs
    type: "mysql"
    host: "db"
    port: 3306
    username: "root"
    password: "secret"
    database: "app"
    optional: false      # If true, failure does not cause non-zero exit
    timeout: 60          # Per-check override (seconds)

  - name: "cache"
    type: "redis"
    host: "redis"
    port: 6379
    optional: true
```

### Per-Check Fields

Use the same fields as the CLI flags (snake_case in YAML):

- Common: `host`, `port`, `username`, `password`, `password_file`, `database`
- DB-specific: `ssl_mode`, `operation_mode`
- HTTP: `http_status_code`, `http_text`, `http_follow_redirect`
- Kafka: `kafka_connection_type`, `kafka_connect_to_leader_via_non_leader`
- DNS: `dns_type`, `dns_expect`, `dns_server`
- InfluxDB: `influx_token`, `influx_org`, `influx_bucket`

See the individual service documentation pages for details and validation rules.

## Execution Behavior

- Checks run sequentially in the order defined.
- Each check uses its own timeout (per-check > global > 30s default).
- **Required checks** ( `optional: false` or omitted): Failure causes non-zero exit.
- **Optional checks** (`optional: true`): Failure prints a warning and continues. Overall exit is 0 if all required checks succeed.
- `fail_fast: true` (default): Stop on the first required failure.
- `fail_fast: false`: Attempt all checks; report the first required failure.

### Output Examples

**All required succeed:**

```
Waiting for primary-db (mysql)...
Wait4it...
Waiting for cache (redis, optional)...
Wait4it.
Success!
```

**Optional fails:**

```
Waiting for primary-db (mysql)...
Wait4it...
Waiting for cache (redis, optional)...
Wait4it.
Warning: optional check "cache" (redis) failed: check failed: context deadline exceeded
Success with 1 optional warning(s)!
```

## Global CLI / Env Options (still apply)

- `--no-banner` / `W4IT_NO_BANNER`
- The `--config` flag and `W4IT_CONFIG` / `W4IT_CONFIG_YAML` vars.

## Notes

- `password_file` is supported per-check and validated early (similar to single-check mode).
- When using a config, single-check flags like `-type`, `-h`, etc. are ignored.
- Unknown check types or validation errors in specs produce clear errors.
- YAML comments and standard features are supported.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0    | All required checks succeeded (optional failures produce warnings only). |
| Non-zero | At least one required check failed (or config parse error). |

See individual service pages for service-specific details. The multi-check feature reuses the same per-service check implementations.