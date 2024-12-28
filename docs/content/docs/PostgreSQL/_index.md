+++
date = '2024-12-24T21:25:42+01:00'
draft = false
title = 'PostgreSQL Check'
+++
PostgreSQL Check validates connectivity to a PostgreSQL database and optionally checks for SSL mode and database availability.

## Usage with Binary
```bash
./wait4it -type=postgres -h=127.0.0.1 -p=5432 -t=60 -u=postgres -P=secret -ssl=disable
```

## Usage with Docker
```bash
docker run ph4r5h4d/wait4it -type=postgres -h=127.0.0.1 -p=5432 -t=60 -u=postgres -P=secret -ssl=disable
```

## Environment Variables

| Variable        | Description                                                               | Default   |
|-----------------|---------------------------------------------------------------------------|-----------|
| W4IT_TYPE       | The type of check (set to `postgres` for PostgreSQL check).               | -         |
| W4IT_TIMEOUT    | Timeout in seconds.                                                      | 30        |
| W4IT_HOST       | The host to check.                                                       | 127.0.0.1 |
| W4IT_PORT       | The port to check on the PostgreSQL host.                                 | 5432      |
| W4IT_USERNAME   | The username for PostgreSQL authentication.                               | -         |
| W4IT_PASSWORD   | The password for PostgreSQL authentication.                               | -         |
| W4IT_DBNAME     | The name of the database to check.                                        | -         |
| W4IT_SSL_MODE   | Whether to enable or disable SSL mode for PostgreSQL (`enable` or `disable`). | `disable` |

## Command-Line Arguments

| Argument        | Description                                                               | Default   |
|-----------------|---------------------------------------------------------------------------|-----------|
| -type           | The type of check (set to `postgres`).                                    | -         |
| -t              | Timeout in seconds.                                                      | 30        |
| -h              | The host to check.                                                       | 127.0.0.1 |
| -p              | The port to check on the PostgreSQL host.                                 | 5432      |
| -u              | The username for PostgreSQL authentication.                               | -         |
| -P              | The password for PostgreSQL authentication.                               | -         |
| -n              | The name of the database to check.                                        | -         |
| -ssl            | Whether to enable or disable SSL mode for PostgreSQL (`disable` or `enable`). | `disable` |

## Notes
{{< callout type="info" >}}
- Environment variables override command-line arguments.
- If `W4IT_SSL_MODE` is not defined, SSL is disabled by default.
- Ensure that the specified PostgreSQL database, username, and password are correct for the connection to succeed.
- {{< /callout >}}

## Exit Codes
| Code | Meaning                            |
|------|------------------------------------|
| 0    | Connection successful.            |
| 1    | Timed out.                         |
| 2    | Validation error or incorrect input.
