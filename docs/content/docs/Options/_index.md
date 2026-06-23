+++
date = '2026-06-23T00:00:00+01:00'
draft = false
title = 'Options'
weight = 20
+++
# Global Options

In addition to the per-service flags documented in the sidebar, Wait4it exposes a few global options that apply to every check. Like all flags, each one can be set either on the command line or via an environment variable (the environment variable is used when the flag is omitted).

## Disable the Startup Banner (`--no-banner`)

By default Wait4it prints an ASCII banner plus sponsor/support links before running a check. If you embed Wait4it in a script, CI pipeline, or a quiet startup sequence, you can suppress this output.

```bash
# Command-line flag
wait4it --no-banner -type=tcp -h=127.0.0.1 -p=8080

# Environment variable
W4IT_NO_BANNER=true wait4it -type=tcp -h=127.0.0.1 -p=8080
```

| Flag        | Environment Variable | Description                          | Default |
|-------------|----------------------|--------------------------------------|---------|
| `--no-banner` | `W4IT_NO_BANNER`     | Disable the startup banner output.   | `false` |

## Timeout (`-t`)

The maximum time Wait4it waits for the target service to become ready, in seconds.

| Flag | Environment Variable | Description                | Default |
|------|----------------------|----------------------------|---------|
| `-t` | `W4IT_TIMEOUT`       | Timeout in seconds.        | 30      |

## Notes
- Environment variables are read only when the corresponding flag is **not** passed on the command line.
- See the per-service pages in the sidebar for the flags specific to each check type.