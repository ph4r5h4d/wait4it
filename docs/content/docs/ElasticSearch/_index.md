+++
date = '2024-12-24T21:25:42+01:00'
draft = false
title = 'ElasitcSearch'
+++
ElasticSearch Check validates connectivity to an ElasticSearch instance and optionally checks for username and password authentication (if the X-Pack extension is enabled).

## Usage with Binary
```bash
./wait4it -type=elasticsearch -h=http://127.0.0.1 -p=9200 -t=60
```

## Usage with Docker
```bash
docker run ph4r5h4d/wait4it -type=elasticsearch -h=http://127.0.0.1 -p=9200 -t=60
```

## Environment Variables

| Variable         | Description                                                               | Default   |
|------------------|---------------------------------------------------------------------------|-----------|
| W4IT_TYPE        | The type of check (set to `elasticsearch` for ElasticSearch check).       | -         |
| W4IT_TIMEOUT     | Timeout in seconds.                                                      | 30        |
| W4IT_HOST        | The host to check.                                                       | 127.0.0.1 |
| W4IT_PORT        | The port to check on the ElasticSearch host.                               | 9200      |

## Command-Line Arguments

| Argument         | Description                                                               | Default   |
|------------------|---------------------------------------------------------------------------|-----------|
| -type            | The type of check (set to `elasticsearch`).                              | -         |
| -t               | Timeout in seconds.                                                      | 30        |
| -h               | The host to check.                                                       | 127.0.0.1 |
| -p               | The port to check on the ElasticSearch host.                              | 9200      |

## Notes
{{< callout type="info" >}}
- Environment variables override command-line arguments.
- Multiple hosts and cluster checks are not supported at the moment.
- Username/password authentication is only supported when the X-Pack extension is enabled. If X-Pack is not activated, the username and password are ignored.
{{< /callout >}}

## Exit Codes
| Code | Meaning                            |
|------|------------------------------------|
| 0    | Connection successful.            |
| 1    | Timed out.                         |
| 2    | Validation error or incorrect input.
