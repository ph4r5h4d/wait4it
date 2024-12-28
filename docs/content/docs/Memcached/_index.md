+++
date = '2024-12-24T21:25:42+01:00'
draft = true
title = 'Memcached'
+++

Memcached Check validates connectivity to a Memcached instance.

## Usage with Binary
./wait4it -type=memcached -h=127.0.0.1 -p=11211 -t=60

## Usage with Docker
docker run ph4r5h4d/wait4it -type=memcached -h=127.0.0.1 -p=11211 -t=60

## Environment Variables

| Variable         | Description                                                               | Default   |
|------------------|---------------------------------------------------------------------------|-----------|
| W4IT_TYPE        | The type of check (set to `memcached` for Memcached check).               | -         |
| W4IT_TIMEOUT     | Timeout in seconds.                                                      | 30        |
| W4IT_HOST        | The host to check.                                                       | 127.0.0.1 |
| W4IT_PORT        | The port to check on the Memcached host.                                  | 11211     |

## Command-Line Arguments

| Argument         | Description                                                               | Default   |
|------------------|---------------------------------------------------------------------------|-----------|
| -type            | The type of check (set to `memcached`).                                   | -         |
| -t               | Timeout in seconds.                                                      | 30        |
| -h               | The host to check.                                                       | 127.0.0.1 |
| -p               | The port to check on the Memcached host.                                  | 11211     |

## Notes
{{< callout type="info" >}}
- Environment variables override command-line arguments.
- Currently, multiple hosts and cluster checks are not supported.
{{< /callout >}}

## Exit Codes
| Code | Meaning                            |
|------|------------------------------------|
| 0    | Connection successful.            |
| 1    | Timed out.                         |
| 2    | Validation error or incorrect input.
