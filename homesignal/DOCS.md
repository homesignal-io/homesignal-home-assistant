# HomeSignal Add-On Notes

## Runtime

The agent is a single Go HTTP server listening on port `8099` by default. Home Assistant ingress is wired to `/ui`.

The container does not request privileged mode, host networking, Docker access, full access, or broad host filesystem mounts. It currently runs the agent as the image default user because Home Assistant owns the mounted `/config` add-on storage path, and the agent must be able to create `/config/device.json` on first boot.

## Identity

On startup, the agent ensures `/config/device.json` exists. If the file is missing, it writes a generated UUIDv4-style `installation_id`. If the file exists, the existing ID is reused.

The add-on uses only the `addon_config:rw` mapping for persistent add-on-owned files. There is no fallback to broad Home Assistant config mounts.

## Options

The agent attempts to read `/data/options.json`. A missing file is accepted and treated as empty configuration. Invalid JSON is an initialization error because it means Supervisor provided malformed options.

## Supervisor And Core API

The add-on requests `hassio_api` and `homeassistant_api` permissions. Home Assistant Supervisor injects `SUPERVISOR_TOKEN` when those APIs are available.

Feature 1 only detects whether the token is present and prepares a placeholder Core API client for:

```text
http://supervisor/core/api/
```

The token is never displayed, persisted, or required for local boot. Missing token produces degraded readiness so local tests and development remain possible.

## Readiness

`/healthz` reports process liveness.

`/readyz` reports initialized local state and Supervisor/Core API availability. Missing `SUPERVISOR_TOKEN` returns HTTP 200 with `degraded: true`; local storage or identity initialization failures prevent startup.
