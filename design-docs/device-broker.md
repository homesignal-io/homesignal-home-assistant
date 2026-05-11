HomeSignal Edge Device Management Requirements

1. Core Architecture

System must use an outbound-only agent model.

Home Assistant host runs HomeSignal Manager add-on.

Agent communicates with HomeSignal cloud over HTTPS/WSS.

No inbound ports required.

No customer LAN access required.

Remote access is optional and separately configured.

Postgres is the canonical source of truth.

WebSocket gateway is transport only, not authority.


2. Device Identity

Each installation must have a stable HomeSignal device ID.

Device identity must survive agent restarts.

Device identity must not depend only on hostname, IP, MAC address, or HA instance name.

Device must generate local installation identity at first run.

Cloud must bind device identity to:
integrator account
customer
site
Home Assistant instance
agent installation


3. Enrollment / Claiming

Unclaimed devices must not receive operational commands.

Claiming requires:
valid provisioning session
logged-in authorized integrator
short-lived pairing code
device-reported unclaimed state

Provisioning sessions must expire.

Pairing codes must be single-use.

Backend must prevent already-claimed devices from being silently re-claimed.

Claiming must create durable credentials.

Temporary claim credentials must be discarded after enrollment.

Claim flow must support stock Home Assistant installs without preloaded hardware.


4. Release / Transfer

System must support release of a device from a site.

Release must revoke cloud credentials.

Release must not break Home Assistant itself.

Released device returns to unclaimed state.

Offline release must still revoke cloud access.

Transfer must be separate from release.

Transfer preserves history while changing controlling account/site ownership.


5. Device Twin / Desired vs Reported State

System must maintain desired state and reported state.

Desired state examples:
required agent version
backup policy
telemetry interval
enabled checks
update policy
managed add-ons
remote access provider metadata

Reported state examples:
HA version
Supervisor version
agent version
hardware type
storage status
backup status
update status
last heartbeat
integration health
remote access configured/not configured

System must compute drift:
desired != reported

Drift must appear in dashboard.

Commands should be generated from desired-state changes where appropriate.

Device twin must be stored in HomeSignal database, not only in memory.


6. Heartbeat / Presence

Agent must send periodic heartbeat.

Heartbeat includes:
device ID
agent version
HA version summary
health summary
timestamp
connection/session ID

Backend must track:
last_seen_at
online/offline
degraded
stale
released/revoked

Online state must be derived from recent heartbeat/socket state.

System must tolerate missed heartbeats without immediate false alarms.

Reconnects must use backoff and jitter.


7. Command Lifecycle

Commands must have durable lifecycle states:
queued
delivered
accepted
running
succeeded
failed
expired
cancelled

Every command must have:
command_id
device_id
requested_by
created_at
expires_at
type
payload
status
result
audit metadata

Agent must ACK command receipt.

Agent must report command result.

Commands must be idempotent where possible.

Commands must be scoped and allowlisted.

Agent must reject unknown command types.

Dangerous commands require stronger confirmation.


8. Supported MVP Commands

request diagnostics
trigger backup
restart HomeSignal add-on
refresh state
check updates
report installed add-ons
report backup status
release local credentials
test remote access metadata

Later:
apply template
stage update
execute update
restore backup
install managed add-on


9. Security Posture

All device-cloud communication must use TLS.

Device credentials must be high-entropy, per-device, scoped, and revocable.

Device tokens must never be exposed in UI.

Device tokens must not be stored in browser localStorage.

Claim codes must not become permanent credentials.

Backend must authorize every command against:
user
integrator account
site
device
role
device state

Agent must enforce local allowlist of permitted operations.

No arbitrary shell execution in MVP.

No unrestricted file read/upload.

No LAN scanning by default.

No subnet routing by default.

No required VPN.

All sensitive actions must be audit logged.


10. Credential Lifecycle

Enrollment issues durable device credential.

Credential can be revoked immediately.

Credential rotation must be supported.

Device must detect revoked credential and enter safe state.

Lost credential recovery requires re-pairing.

Future option:
device-generated keypair + signed device certificate / mTLS.


11. Access Control

Roles:
Owner
Admin
Technician
Read-only

Permissions must cover:
create site
claim device
release device
transfer device
view health
run commands
manage backups
manage team
configure remote access
view audit log

Technicians should only access assigned sites where possible.


12. Remote Access

Remote access is optional.

HomeSignal must support storing remote access metadata:
provider
URL
node name
notes
last verified
configured_by

Supported providers:
manual URL
Tailscale
Nabu Casa
Cloudflare Tunnel
VPN/custom

Remote access link must be separate from agent health.

Browser reachability check is convenience only, not authoritative health.

HomeSignal must not require access to customer private networks.


13. Backup Management

Agent must report backup status.

System must track:
last successful backup
last failed backup
backup location metadata
backup policy
retention policy
backup size
failure reason

MVP may only monitor/trigger local HA backups.

Later versions may support encrypted offsite backup.


14. Update Management

Agent must report:
HA Core update availability
HAOS update availability
Supervisor update availability
agent update availability
managed add-on update availability

System must support update policy:
manual
notify only
approved window
blocked version
staged rollout

MVP should not auto-update HA without explicit approval.


15. Diagnostics

Agent must collect safe diagnostic bundle:
versions
health summary
logs from HomeSignal add-on
Supervisor status
recent command results
backup status
storage summary

Diagnostics upload must be explicit and audited.

Sensitive secrets must be redacted.

Large diagnostics must use HTTPS upload, not WebSocket.


16. Audit Logging

Audit log must capture:
login
site creation
device claim
device release
device transfer
command requested
command completed
credential revoked
remote URL changed
team/user permission changes

Audit entries include:
actor
timestamp
account
site
device
action
result
source IP if available


17. Failure Handling

If cloud is unreachable:
agent continues local HA operation
agent queues non-dangerous reports locally within limit
agent reconnects with backoff

If agent is offline:
cloud marks device offline/stale
commands remain queued until expiry

If device is released while offline:
cloud revokes credential immediately
device is denied on next reconnect

If duplicate device identity appears:
backend must block or quarantine session

If pairing race occurs:
first valid claim wins
all others fail


18. Data Model Minimum

accounts
users
roles
customers
sites
devices
device_credentials
provisioning_sessions
device_reported_state
device_desired_state
commands
command_results
heartbeats
alerts
backups
remote_access_links
audit_events


19. Alerts

Initial alert types:
device offline
heartbeat stale
backup failed
backup overdue
agent version outdated
HA update available
storage high
credential revoked
claim failed
command failed

Alerts must support:
severity
acknowledge
resolve
snooze later


20. Non-Goals for MVP

No custom VPN.

No subnet routing.

No arbitrary remote shell.

No full HA fork.

No custom HAOS image required.

No mandatory Tailscale.

No AWS IoT Core unless architecture changes.

No high-frequency telemetry pipeline.

No customer LAN inventory scanning.


21. MVP Success Criteria

Integrator can:
create site
install HomeSignal add-on
claim HA device
see online/offline status
see basic HA health
trigger diagnostics
trigger backup
see backup status
store remote access URL
release device safely

System can:
securely enroll device
prevent stale/replayed claims
revoke device access
track desired/reported state
durably queue commands
audit all sensitive operations
operate without inbound network access
