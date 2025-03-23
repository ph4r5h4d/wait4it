+++
date = '2025-03-23T11:25:42+01:00'
draft = false
title = 'Kafka'
+++
Kafka Check validates connectivity to a Kafka instance. This can check the connection to the leader instance,
or try to connect to the leader via a non-leader instance.  
**This module support is beta. No authentication scheme is currently supported**

## Usage with Binary
```bash
./wait4it -type=kafka -h=127.0.0.1 -p=9092 -kafka-connect-to-leader-via-non-leader=false -t=60
```

## Usage with Docker
```bash
docker run ph4r5h4d/wait4it -type=kafka -h=127.0.0.1 -p=9092 -kafka-connect-to-leader-via-non-leader=false -t=60
```

## Environment Variables
| Variable                             | Description                                            | Default   |
|--------------------------------------|--------------------------------------------------------|-----------|
| W4IT_TYPE                            | The type of check (set to `kafka` for Kafka check).    | -         |
| W4IT_TIMEOUT                         | Timeout in seconds.                                    | 30        |
| W4IT_HOST                            | The host to check.                                     | 127.0.0.1 |
| W4IT_PORT                            | The port to check on the Kafka host.                   | 9092      |
| W4IT_KAFKA_CONNECTION_TYPE           | The connection type toward Kafka                       | tcp       |
| W4IT_KAFKA_CONNECT_TO_LEADER_VIA_NON_LEADER | Whether to try to connect to a leader via a non-leader | false     |

## Command-Line Arguments
| Argument                             | Description                                                               | Default   |
|--------------------------------------|---------------------------------------------------------------------------|-----------|
| -type                                | The type of check (set to `kafka` for Kafka check).                       | -         |
| -t                                   | Timeout in seconds.                                                       | 30        |
| -h                                   | The host to check.                                                        | 127.0.0.1 |
| -p                                   | The port to check on the Kafka host.                                      | 9092      |
| -kafka-connection-type               | The connection type toward Kafka                                          | tcp       |
| -kafka-connect-to-leader-via-non-leader | Whether to try to connect to a leader via a non-leader                     | false     |
## Notes

{{< callout type="info" >}}
- Environment variables override command-line arguments.
- Authentication scheme is not supported yet.
- If you have a requirement for authentication, please open an issue.
{{< /callout >}}

## Exit Codes
| Code | Meaning                            |
|------|------------------------------------|
| 0    | Connection successful.            |
| 1    | Timed out.                         |
| 2    | Validation error or incorrect input.|
