+++
date = '2024-12-24T21:25:42+01:00'
draft = false
title = 'MongoDB'
+++
MongoDB Check validates connectivity to a MongoDB instance.

## Usage with Binary
```bash
./wait4it -type=mongo -h=127.0.0.1 -p=27017 -t=60 -u=mongoadmin -P=secret
```

## Usage with Docker
```bash
docker run ph4r5h4d/wait4it -type=mongo -h=127.0.0.1 -p=27017 -t=60 -u=mongoadmin -P=secret
```

## Environment Variables

| Variable           | Description                                           | Default   |
|--------------------|-------------------------------------------------------|-----------|
| W4IT_TYPE          | The type of check (set to `mongo` for MongoDB check). | -         |
| W4IT_TIMEOUT       | Timeout in seconds.                                   | 30        |
| W4IT_HOST          | The host to check.                                    | 127.0.0.1 |
| W4IT_PORT          | The port to check on the MongoDB host.                | 27017     |
| W4IT_USERNAME      | The username for MongoDB authentication.              | -         |
| W4IT_PASSWORD      | The password for MongoDB authentication.              | -         |
| W4IT_PASSWORD_FILE | The file with password for MongoDB authentication.    | -         |
## Command-Line Arguments

| Argument | Description                                        | Default   |
|----------|----------------------------------------------------|-----------|
| -type    | The type of check (set to `mongo`).                | -         |
| -t       | Timeout in seconds.                                | 30        |
| -h       | The host to check.                                 | 127.0.0.1 |
| -p       | The port to check on the MongoDB host.             | 27017     |
| -u       | The username for MongoDB authentication.           | -         |
| -P       | The password for MongoDB authentication.           | -         |
| -Pf      | The file with password for MongoDB authentication. | -         |

## Notes
{{< callout type="info" >}}
- Environment variables override command-line arguments.
- Currently, only username/password authentication mechanism is supported.
{{< /callout >}}

## Exit Codes
| Code | Meaning                              |
|------|--------------------------------------|
| 0    | Connection successful.               |
| 1    | Timed out.                           |
| 2    | Validation error or incorrect input. |
