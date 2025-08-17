+++
date = '2024-12-24T21:25:42+01:00'
draft = false
title = 'Redis'
+++
Redis Check validates connectivity to a Redis instance and optionally checks for authentication and operation mode (standalone or cluster).

## Usage with Binary
```bash
./wait4it -type=redis -h=127.0.0.1 -p=6379 -t=60 -P=secret
```

## Usage with Docker
```bash
docker run ph4r5h4d/wait4it -type=redis -h=127.0.0.1 -p=6379 -t=60 -P=secret
```

## Environment Variables

| Variable            | Description                                               | Default      |
|---------------------|-----------------------------------------------------------|--------------|
| W4IT_TYPE           | The type of check (set to `redis` for Redis check).       | -            |
| W4IT_TIMEOUT        | Timeout in seconds.                                       | 30           |
| W4IT_HOST           | The host to check.                                        | 127.0.0.1    |
| W4IT_PORT           | The port to check on the Redis host.                      | 6379         |
| W4IT_PASSWORD       | The password for Redis authentication (if any).           | -            |
| W4IT_PASSWORD_FILE  | The file with password for Redis authentication (if any). | -            |
| W4IT_OPERATION_MODE | The Redis operation mode (`standalone` or `cluster`).     | `standalone` |

## Command-Line Arguments

| Argument        | Description                                               | Default      |
|-----------------|-----------------------------------------------------------|--------------|
| -type           | The type of check (set to `redis`).                       | -            |
| -t              | Timeout in seconds.                                       | 30           |
| -h              | The host to check.                                        | 127.0.0.1    |
| -p              | The port to check on the Redis host.                      | 6379         |
| -P              | The password for Redis authentication (if any).           | -            |
| -Pf             | The file with password for Redis authentication (if any). | -            |
| -operation-mode | The Redis operation mode (`standalone` or `cluster`).     | `standalone` |

## Notes
{{< callout type="info" >}}
- Environment variables override command-line arguments.
- If `W4IT_OPERATION_MODE` is not defined, it defaults to `standalone`.
- Redis Sentinel is not supported yet.
- This version can only check one host within a Redis cluster; support for checking multiple hosts in a Redis cluster might be added upon request.
{{< /callout >}}


## Exit Codes
| Code | Meaning                              |
|------|--------------------------------------|
| 0    | Connection successful.               |
| 1    | Timed out.                           |
| 2    | Validation error or incorrect input. |
