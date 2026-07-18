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

