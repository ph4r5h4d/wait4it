+++
date = '2024-12-24T21:25:42+01:00'
draft = false
title = 'HTTP Check'
+++
HTTP Check validates connectivity to an HTTP endpoint and optionally checks for specific status codes or text in the response.

## Usage with Binary
```bash
./wait4it -type=http -h=https://example.com -t=60 -http-status=200 -http-text="Welcome"
```
## Usage with Docker
```bash
docker run ph4r5h4d/wait4it -type=http -h=https://example.com -t=60 -http-status=200 -http-text="Welcome"
```
## Environment Variables

| Variable                   | Description                                                                                     | Default   |
|----------------------------|-------------------------------------------------------------------------------------------------|-----------|
| W4IT_TYPE                  | The type of check (set to `http` for HTTP check).                                              | -         |
| W4IT_TIMEOUT               | Timeout in seconds.                                                                            | 30        |
| W4IT_HOST                  | The HTTP endpoint to check.                                                                    | 127.0.0.1 |
| W4IT_HTTP_STATUS_CODE      | Expected HTTP status code.                                                                     | 200       |
| W4IT_HTTP_TEXT             | Substring to search for in the HTTP response.                                                  | -         |
| W4IT_HTTP_FOLLOW_REDIRECT  | Whether to follow redirects (`true` or `false`).                                               | true      |

## Command-Line Arguments

| Argument                   | Description                                                                                     | Default   |
|----------------------------|-------------------------------------------------------------------------------------------------|-----------|
| -type                      | The type of check (set to `http` for HTTP check).                                              | -         |
| -t                         | Timeout in seconds.                                                                            | 30        |
| -h                         | The HTTP endpoint to check.                                                                    | 127.0.0.1 |
| -http-status               | Expected HTTP status code.                                                                     | 200       |
| -http-text                 | Substring to search for in the HTTP response.                                                  | -         |
| -http-follow-redirect      | Whether to follow redirects (`true` or `false`).                                               | true      |

## Notes
{{< callout type="info" >}}
- Environment variables override command-line arguments.
- If `W4IT_HTTP_STATUS_CODE` is not defined, the check defaults to expecting status code `200`.
- If `W4IT_HTTP_TEXT` is not defined, only the status code is checked.
{{< /callout >}}

## Exit Codes
| Code | Meaning                            |
|------|------------------------------------|
| 0    | Connection successful.            |
| 1    | Timed out.                         |
| 2    | Validation error or incorrect input.

