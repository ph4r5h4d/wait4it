+++
date = '2026-05-27T21:00:00+00:00'
draft = false
title = 'InfluxDB Check'
+++

InfluxDB Check validates connectivity and readiness of an InfluxDB 2.x instance (using the `/health` endpoint and requiring status `pass`).

## Usage with Binary
```bash
./wait4it -type=influxdb -h=127.0.0.1 -p=8086 -t=60 -token=my-token
```

## Usage with Docker
```bash
docker run ph4r5h4d/wait4it -type=influxdb -h=127.0.0.1 -p=8086 -t=60 -token=my-token
```

## Environment Variables

| Variable            | Description                                                                 | Default |
|---------------------|-----------------------------------------------------------------------------|---------|
| W4IT_TYPE           | The type of check (set to `influxdb`).                                      | -       |
| W4IT_TIMEOUT        | Timeout in seconds.                                                         | 30      |
| W4IT_HOST           | The host to check.                                                          | 127.0.0.1 |
| W4IT_PORT           | The port to check (InfluxDB default is 8086).                               | 8086    |
| W4IT_INFLUX_TOKEN   | API token for InfluxDB (optional for basic health checks).                  | -       |
| W4IT_INFLUX_ORG     | Organization name (optional; not used by current basic health check).       | -       |
| W4IT_INFLUX_BUCKET  | Bucket name (optional; not used by current basic health check).             | -       |

## Command-Line Arguments

| Argument | Description                                                              | Default |
|----------|--------------------------------------------------------------------------|---------|
| -type    | The type of check (set to `influxdb`).                                   | -       |
| -t       | Timeout in seconds.                                                      | 30      |
| -h       | The host to check.                                                       | 127.0.0.1 |
| -p       | The port to check.                                                       | 8086    |
| -token   | InfluxDB API token (optional for basic health checks).                   | -       |
| -org     | Organization name (optional; currently unused).                          | -       |
| -bucket  | Bucket name (optional; currently unused).                                | -       |

## Notes
{{< callout type="info" >}}
- The `/health` endpoint in InfluxDB 2.x is publicly accessible and does not require a token in standard OSS configurations. The checker succeeds only when the reported status is `pass`.
- Token/Org/Bucket flags are accepted (for forward compatibility) but the current implementation only performs a basic instance health check and does not use Org/Bucket.
- For full functionality (writes/queries) a valid token with appropriate permissions is required.
- This checker targets InfluxDB 2.x. InfluxDB 1.x users can often use the generic `tcp` or `http` checkers. InfluxDB 3.x support may be added in a future release.
- HTTPS: pass scheme via `-h` (e.g. `-h=https://example -p=8086`).
- {{< /callout >}}

## Exit Codes
| Code | Meaning                              |
|------|--------------------------------------|
| 0    | Connection / readiness successful.   |
| 1    | Timed out.                           |
| 2    | Validation error or incorrect input. |
