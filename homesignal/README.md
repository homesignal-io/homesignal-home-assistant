# HomeSignal

HomeSignal is a local-build Home Assistant add-on skeleton for the HomeSignal agent. This first version only provides local identity, status endpoints, an ingress placeholder UI, and Supervisor/Core API permission wiring.

## Install Locally

1. Copy or clone this repository into a Home Assistant add-on repository location.
2. In Home Assistant, go to **Settings > Add-ons > Add-on Store**.
3. Add the local repository path or refresh the local add-on repository.
4. Install the **HomeSignal** add-on.
5. Start the add-on and open its Web UI.

This skeleton intentionally omits a production `image` field in `config.yaml`, so Home Assistant can build it from this add-on folder.

## Endpoints

- `/healthz`: process liveness.
- `/readyz`: identity/config readiness and degraded Supervisor/Core API status.
- `/version`: build metadata.
- `/ui`: basic Home Assistant ingress placeholder page.

## Permissions

The add-on requests only:

- `hassio_api: true`
- `homeassistant_api: true`
- `addon_config:rw`

It does not request Docker access, host networking, privileged mode, full access, the Docker socket, or broad Home Assistant filesystem mappings.

## Storage

The agent stores add-on-owned data in `/config`, backed by Home Assistant's `addon_config:rw` mapping. On first boot it creates:

```text
/config/device.json
```

The file contains a generated `installation_id` and is reused across restarts.

## Current Limitations

This release does not implement enrollment, cloud authentication, telemetry, topology discovery, backup actions, update orchestration, IoT Core, or command execution.
