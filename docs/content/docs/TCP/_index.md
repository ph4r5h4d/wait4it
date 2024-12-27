+++
date = '2024-12-24T21:25:42+01:00'
draft = false
title = 'TCP Port Check'
+++
TCP Port Check validates connectivity to a specific TCP port on a target host within a given timeout.

## Usage with Binary
```bash
./wait4it -type=tcp -h=127.0.0.1 -p=8080 -t=60
```

## Usage with Docker
```bash
docker run ph4r5h4d/wait4it -type=tcp -h=127.0.0.1 -p=8080 -t=60
```

## Environment Variables

| Variable       | Description                           | Default   |
|----------------|---------------------------------------|-----------|
| W4IT_TYPE      | The type of check (set to `tcp`).     | -         |
| W4IT_TIMEOUT   | Timeout in seconds.                  | 30        |
| W4IT_HOST      | The host to check.                   | 127.0.0.1 |
| W4IT_PORT      | The TCP port to check.               | -         |

## Command-Line Arguments

| Argument | Description                           | Default   |
|----------|---------------------------------------|-----------|
| -type    | The type of check (set to `tcp`).     | -         |
| -t       | Timeout in seconds.                  | 30        |
| -h       | The host to check.                   | 127.0.0.1 |
| -p       | The TCP port to check.               | -         |

## Notes
- Environment variables override command-line arguments.
- Ensure the target host and port are reachable within the timeout period.

## Exit Codes
| Code | Meaning                            |
|------|------------------------------------|
| 0    | Connection successful.            |
| 1    | Timed out.                         |
| 2    | Validation error or incorrect input.
