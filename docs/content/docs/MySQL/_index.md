+++
date = '2024-12-24T21:25:42+01:00'
draft = false
title = 'MySQL Check'
+++
MySQL Check validates connectivity to a MySQL database and optionally checks for specific database availability.

## Usage with Binary
```bash
./wait4it -type=mysql -h=127.0.0.1 -p=3306 -t=60 -u=root -P=secret -n=app
```

## Usage with Docker
```bash
docker run ph4r5h4d/wait4it -type=mysql -h=127.0.0.1 -p=3306 -t=60 -u=root -P=secret -n=app
```

## Environment Variables

| Variable           | Description                                         | Default   |
|--------------------|-----------------------------------------------------|-----------|
| W4IT_TYPE          | The type of check (set to `mysql` for MySQL check). | -         |
| W4IT_TIMEOUT       | Timeout in seconds.                                 | 30        |
| W4IT_HOST          | The host to check.                                  | 127.0.0.1 |
| W4IT_PORT          | The port to check on the MySQL host.                | 3306      |
| W4IT_USERNAME      | The username for MySQL authentication.              | -         |
| W4IT_PASSWORD      | The password for MySQL authentication.              | -         |
| W4IT_PASSWORD_FILE | The file with password for MySQL authentication.    | -         |
| W4IT_DBNAME        | The name of the database to check.                  | -         |

## Command-Line Arguments

| Argument | Description                                      | Default   |
|----------|--------------------------------------------------|-----------|
| -type    | The type of check (set to `mysql`).              | -         |
| -t       | Timeout in seconds.                              | 30        |
| -h       | The host to check.                               | 127.0.0.1 |
| -p       | The port to check on the MySQL host.             | 3306      |
| -u       | The username for MySQL authentication.           | -         |
| -P       | The password for MySQL authentication.           | -         |
| -Pf      | The file with password for MySQL authentication. | -         |
| -n       | The name of the database to check.               | -         |

## Notes
{{< callout type="info" >}}
- Environment variables override command-line arguments.
- Ensure that the specified MySQL database and username/password are correct for the connection to succeed.
{{< /callout >}}

## Exit Codes
| Code | Meaning                              |
|------|--------------------------------------|
| 0    | Connection successful.               |
| 1    | Timed out.                           |
| 2    | Validation error or incorrect input. |
