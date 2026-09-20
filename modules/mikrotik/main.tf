terraform {
  required_providers {
    routeros = {
      source = "terraform-routeros/routeros"
    }
  }
}

variable "ipv6_bridge_ip" {
  type        = string
  description = "IPv6 address for the bridge interface"
}

variable "ipv6_prefix" {
  type        = string
  description = "Prefix length for the IPv6 address"
}

resource "routeros_system_logging_action" "loki" {
  name = "loki"
  target = "remote"
  remote = "192.168.1.58"
  remote_port = 10514
  remote_protocol = "udp"
  remote_log_format = "syslog"
  syslog_facility = "local0"
  syslog_severity = "info"
}

resource "routeros_system_logging" "loki" {
  action = routeros_system_logging_action.loki.name
  topics = ["info"]
}

