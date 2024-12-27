+++
date = '2024-12-24T21:25:42+01:00'
draft = false
title = 'Aerospike'
+++
Aerospike Check validates connectivity to an Aerospike instance.

## Usage with Binary
```bash
./wait4it -type=aerospike -h=127.0.0.1 -p=3000 -t=60
```

## Usage with Docker
```bash
docker run ph4r5h4d/wait4it -type=aerospike -h=127.0.0.1 -p=3000 -t=60
```

## Environment Variables

| Variable         | Description                                                               | Default   |
|------------------|---------------------------------------------------------------------------|-----------|
| W4IT_TYPE        | The type of check (set to `aerospike` for Aerospike check).               | -         |
| W4IT_TIMEOUT     | Timeout in seconds.                                                      | 30        |
| W4IT_HOST        | The host to check.                                                       | 127.0.0.1 |
| W4IT_PORT        | The port to check on the Aerospike host.                                  | 3000      |

## Command-Line Arguments

| Argument         | Description                                                               | Default   |
|------------------|---------------------------------------------------------------------------|-----------|
| -type            | The type of check (set to `aerospike`).                                  | -         |
| -t               | Timeout in seconds.                                                      | 30        |
| -h               | The host to check.                                                       | 127.0.0.1 |
| -p               | The port to check on the Aerospike host.                                  | 3000      |

## Notes
{{< callout type="info" >}}
- Environment variables override command-line arguments.
- Multiple hosts and cluster checks are not supported at the moment.
{{< /callout >}}

## Exit Codes
| Code | Meaning                            |
|------|------------------------------------|
| 0    | Connection successful.            |
| 1    | Timed out.                         |
| 2    | Validation error or incorrect input.
