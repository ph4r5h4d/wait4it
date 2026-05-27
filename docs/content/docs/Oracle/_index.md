+++
date = '2026-05-27T12:00:00+00:00'
draft = false
title = 'Oracle Check'
+++

Oracle Check validates connectivity to an Oracle Database (21c and higher, including 23ai) using a service name (typically a PDB such as `FREEPDB1`).

## Usage with Binary
```bash
./wait4it -type=oracle -h=127.0.0.1 -p=1521 -t=60 -u=app -P=secret -n=FREEPDB1
```

## Usage with Docker
```bash
docker run ph4r5h4d/wait4it -type=oracle -h=127.0.0.1 -p=1521 -t=60 -u=app -P=secret -n=FREEPDB1
```

## Environment Variables

| Variable           | Description                                                              | Default   |
|--------------------|--------------------------------------------------------------------------|-----------|
| W4IT_TYPE          | The type of check (set to `oracle` for Oracle check).                    | -         |
| W4IT_TIMEOUT       | Timeout in seconds.                                                      | 30        |
| W4IT_HOST          | The host to check.                                                       | 127.0.0.1 |
| W4IT_PORT          | The port to check on the Oracle host (default Oracle listener port).     | 1521      |
| W4IT_USERNAME      | The username for Oracle authentication.                                  | -         |
| W4IT_PASSWORD      | The password for Oracle authentication.                                  | -         |
| W4IT_PASSWORD_FILE | The file with password for Oracle authentication.                        | -         |
| W4IT_DBNAME        | The service name (or PDB name, e.g. `FREEPDB1` or `XEPDB1`) to connect to. | -         |

## Command-Line Arguments

| Argument | Description                                                              | Default   |
|----------|--------------------------------------------------------------------------|-----------|
| -type    | The type of check (set to `oracle`).                                     | -         |
| -t       | Timeout in seconds.                                                      | 30        |
| -h       | The host to check.                                                       | 127.0.0.1 |
| -p       | The port to check on the Oracle host.                                    | 1521      |
| -u       | The username for Oracle authentication.                                  | -         |
| -P       | The password for Oracle authentication.                                  | -         |
| -Pf      | The file with password for Oracle authentication.                        | -         |
| -n       | The service name / PDB name to connect to (recommended for 21c+).        | -         |

## Notes
{{< callout type="info" >}}
- Environment variables override command-line arguments.
- For Oracle 21c and higher (including 23ai Free / XE), use the PDB service name in `-n` / `W4IT_DBNAME` (commonly `FREEPDB1` with popular container images such as `gvenzl/oracle-free` or `gvenzl/oracle-xe`).
- The check performs a TCP connection + authentication + `PingContext`. It does not run any DDL/DML.
- For Oracle Autonomous Database or wallet-based TLS connections, additional configuration (not yet exposed via CLI flags) would be required.
- {{< /callout >}}

## Exit Codes
| Code | Meaning                              |
|------|--------------------------------------|
| 0    | Connection successful.               |
| 1    | Timed out.                           |
| 2    | Validation error or incorrect input. |
