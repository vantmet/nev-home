# Terraform component diagrams

## Graphana module

```mermaid
flowchart LR
    User["User / Browser"] -->|HTTP 3000| Grafana["docker_container.grapaha-main\nGrafana"]
    Grafana -->|Uses image| Image["docker_image.graphana\ngrafana/grafana:latest"]
    Grafana -->|Attached to| Network["docker_network.tig\nBridge network (IPv6)"]
    Grafana -->|Persists config| Volume["docker_volume.grafana_config"]
    Grafana -->|Environment| Env["GRAFANA_ADMIN_PASSWORD"]
```

This module creates a Docker bridge network, pulls the Grafana image, mounts a persistent volume for `/var/lib/graphana`, and runs a container exposed on port `3000`.

## Mikrotik module

```mermaid
flowchart LR
    Router["RouterOS device"] -->|Info logs| Syslog["routeros_system_logging.loki"]
    Syslog -->|Uses action| Action["routeros_system_logging_action.loki"]
    Action -->|UDP 10514| Loki["Remote collector\n192.168.1.58:10514"]
```

This module configures RouterOS to send `info` level syslog entries to a remote Loki endpoint over UDP on port `10514`.
