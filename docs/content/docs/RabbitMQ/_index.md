+++
date = '2024-12-24T21:25:42+01:00'
draft = false
title = 'RabbitMQ'
+++
RabbitMQ Check validates connectivity to a RabbitMQ instance and optionally checks for username and password authentication.

## Usage with Binary
./wait4it -type=rabbitmq -h=127.0.0.1 -p=5267 -t=60 -u=guest -P=guest

## Usage with Docker
docker run ph4r5h4d/wait4it -type=rabbitmq -h=127.0.0.1 -p=5267 -t=60 -u=guest -P=guest

## Environment Variables

| Variable           | Description                                               | Default   |
|--------------------|-----------------------------------------------------------|-----------|
| W4IT_TYPE          | The type of check (set to `rabbitmq` for RabbitMQ check). | -         |
| W4IT_TIMEOUT       | Timeout in seconds.                                       | 30        |
| W4IT_HOST          | The host to check.                                        | 127.0.0.1 |
| W4IT_PORT          | The port to check on the RabbitMQ host.                   | 5267      |
| W4IT_USERNAME      | The username for RabbitMQ authentication.                 | -         |
| W4IT_PASSWORD      | The password for RabbitMQ authentication.                 | -         |
| W4IT_PASSWORD_FILE | The file with password for RabbitMQ authentication.       | -         |

## Command-Line Arguments

| Argument | Description                                         | Default   |
|----------|-----------------------------------------------------|-----------|
| -type    | The type of check (set to `rabbitmq`).              | -         |
| -t       | Timeout in seconds.                                 | 30        |
| -h       | The host to check.                                  | 127.0.0.1 |
| -p       | The port to check on the RabbitMQ host.             | 5267      |
| -u       | The username for RabbitMQ authentication.           | -         |
| -P       | The password for RabbitMQ authentication.           | -         |
| -Pf      | The file with password for RabbitMQ authentication. | -         |

## Notes
{{< callout type="info" >}}
- Environment variables override command-line arguments.
- Ensure that the specified RabbitMQ username and password are correct for the connection to succeed.
- {{< /callout >}}

## Exit Codes
| Code | Meaning                              |
|------|--------------------------------------|
| 0    | Connection successful.               |
| 1    | Timed out.                           |
| 2    | Validation error or incorrect input. |
