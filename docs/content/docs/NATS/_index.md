+++
date = '2026-10-15T00:00:00+01:00'
draft = false
title = 'NATS'
+++
NATS Check validates connectivity to a NATS server by completing a client connection and NATS protocol handshake.

## Usage with Binary
./wait4it -type=nats -h=127.0.0.1 -p=4222 -t=60

With username and password:

./wait4it -type=nats -h=127.0.0.1 -p=4222 -t=60 -u=app -P=secret

With token authentication:

./wait4it -type=nats -h=127.0.0.1 -p=4222 -t=60 -nats-token=MyToken

With TLS (default system trust; no custom CA in v1):

./wait4it -type=nats -h=127.0.0.1 -p=4222 -t=60 -ssl=enable

## Usage with Docker
docker run ph4r5h4d/wait4it -type=nats -h=127.0.0.1 -p=4222 -t=60

## Environment Variables

| Variable           | Description                                      | Default   |
|--------------------|--------------------------------------------------|-----------|
| W4IT_TYPE          | The type of check (set to `nats` for NATS).      | -         |
| W4IT_TIMEOUT       | Timeout in seconds.                              | 30        |
| W4IT_HOST          | The host to check.                               | 127.0.0.1 |
| W4IT_PORT          | The NATS client port (typically 4222).           | -         |
| W4IT_USERNAME      | Username for NATS user/password authentication.  | -         |
| W4IT_PASSWORD      | Password for NATS user/password authentication.  | -         |
| W4IT_PASSWORD_FILE | File containing the NATS password.               | -         |
| W4IT_NATS_TOKEN    | Token for NATS token authentication.             | -         |
| W4IT_SSL_MODE      | Set to `enable` or `require` for TLS (`tls://`). | disable   |

## Command-Line Arguments

| Argument      | Description                                      | Default   |
|---------------|--------------------------------------------------|-----------|
| -type         | The type of check (set to `nats`).               | -         |
| -t            | Timeout in seconds.                              | 30        |
| -h            | The host to check.                               | 127.0.0.1 |
| -p            | The NATS client port (typically 4222).           | -         |
| -u            | Username for NATS user/password authentication.  | -         |
| -P            | Password for NATS user/password authentication.  | -         |
| -Pf           | File containing the NATS password.               | -         |
| -nats-token   | Token for NATS token authentication.             | -         |
| -ssl          | Set to `enable` or `require` for TLS.            | disable   |

## Notes
{{< callout type="info" >}}
- Environment variables override command-line arguments when the corresponding flag is omitted.
- You must set `-p` / `W4IT_PORT` (for example `4222`); the global port default is `0`, which is invalid for NATS.
- Token and username/password can both be supplied; the NATS client may use either depending on server configuration.
- JetStream readiness, NKey/credentials files, and custom TLS CA paths are not supported in this version.
{{< /callout >}}

## Exit Codes
| Code | Meaning                              |
|------|--------------------------------------|
| 0    | Connection successful.               |
| 1    | Timed out.                           |
| 2    | Validation error or incorrect input. |
