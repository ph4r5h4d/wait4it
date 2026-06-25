+++
date = '2026-06-24T00:00:00+01:00'
draft = false
title = 'Multi-Check YAML Configuration'
weight = 15
+++
# Multi-Check YAML Configuration

Wait4it supports defining **one or more checks** using a YAML file (via `--config` / `W4IT_CONFIG`) or inline YAML content (via the `W4IT_CONFIG_YAML` environment variable).

A single check can also be expressed in YAML (see example below). When a configuration source is provided, all other check-specific flags and their environment variable equivalents are ignored.

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

### Using Docker

```bash
docker run -v $(pwd)/checks.yaml:/config.yaml ph4r5h4d/wait4it --config /config.yaml
```

## Single Check via YAML

You can use the same YAML format for a **single check**:

```yaml
version: "1"
checks:
  - type: tcp
    host: 127.0.0.1
    port: 8080
    timeout: 10
```

```bash
wait4it --config single-check.yaml
```

This is useful if you want a consistent configuration format even when only one service is involved.

## YAML Schema

```yaml
version: "1"              # Recommended. Only "1" is currently supported.
timeout: 30               # Global default timeout (seconds). Default: 30
fail_fast: true           # Stop on first required failure. Default: true

checks:                   # Required. At least one check.
  - name: "string"        # Optional. Used in log messages.
    type: "string"        # Required. One of the supported check types (tcp, mysql, redis, etc.).
    host: "string"        # Required for most types.
    port: 0               # Required for most types.
    username: "string"
    password: "string"
    password_file: "string"
    database: "string"

    optional: false       # If true, failure is non-fatal (warning + continue). Default: false
    timeout: 0            # Per-check timeout override (seconds).

    # Database-specific
    ssl_mode: "disable"
    operation_mode: "standalone"

    # HTTP-specific
    http_status_code: 200
    http_text: "string"
    http_follow_redirect: true

    # Kafka-specific
    kafka_connection_type: "tcp"
    kafka_connect_to_leader_via_non_leader: false

    # DNS-specific
    dns_type: "A"
    dns_expect: "string"
    dns_server: "string"

    # InfluxDB-specific
    influx_token: "string"
    influx_org: "string"
    influx_bucket: "string"
```

See the individual service documentation pages for the exact meaning and validation rules of each field.

## Execution Behavior

- Checks run **sequentially** (not concurrently) in the order defined.
- Each check uses its own timeout (per-check > global > 30s default).
- **Required checks** (`optional: false` or omitted): Failure causes non-zero exit.
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

## Examples

### 1. Single Check (TCP)

```yaml
version: "1"
checks:
  - type: tcp
    host: 127.0.0.1
    port: 8080
    timeout: 10
```

### 2. Multiple Checks with One Optional

```yaml
version: "1"
timeout: 30
fail_fast: true

checks:
  - name: "primary-db"
    type: mysql
    host: db
    port: 3306
    username: root
    password: secret
    database: app
    optional: false

  - name: "cache"
    type: redis
    host: redis
    port: 6379
    optional: true
```

### 3. Using fail_fast: false

```yaml
version: "1"
fail_fast: false

checks:
  - name: "must-start-first"
    type: tcp
    host: 10.0.0.5
    port: 5432
    timeout: 60

  - name: "optional-later"
    type: http
    host: http://10.0.0.10/health
    http_status_code: 200
    optional: true
    timeout: 15
```

### 4. Mixed Services with Password File

```yaml
version: "1"
timeout: 45

checks:
  - name: "postgres"
    type: postgres
    host: postgres
    port: 5432
    username: appuser
    password_file: /secrets/db-password
    database: myapp
    ssl_mode: require

  - name: "elasticsearch"
    type: elasticsearch
    host: es
    port: 9200
    optional: true
```

### 5. Inline via Environment Variable (Single or Multi)

```bash
W4IT_CONFIG_YAML='version: "1"
checks:
  - type: kafka
    host: kafka
    port: 9092
    timeout: 20' wait4it
```

### 6. Per-Check Timeouts + Optional

```yaml
version: "1"
timeout: 30

checks:
  - type: mysql
    host: primary-db
    port: 3306
    username: root
    password: secret
    timeout: 120          # slow startup DB gets more time

  - type: redis
    host: redis
    port: 6379
    optional: true
    timeout: 5
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