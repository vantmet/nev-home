resource "routeros_ip_vrf" "vrf_main" {
  name       = "main"
  interfaces = ["all"]
}
resource "routeros_ipv6_pool" "BRSKPool" {
  name          = "BRSKPool"
  prefix        = var.ipv6_prefix
  prefix_length = 64
}
resource "routeros_ipv6_address" "lanv6" {
  address   = var.ipv6_bridge_ip
  interface = routeros_interface_bridge.bridge.name
  advertise = true
  from_pool = routeros_ipv6_pool.BRSKPool.name
  no_dad    = true

}
resource "routeros_ip_address" "lan" {
  address   = "192.168.1.1/24"
  interface = routeros_interface_bridge.bridge.name
  network   = "192.168.1.0"
  comment   = "defconf"
  vrf       = routeros_ip_vrf.vrf_main.name
}

resource "routeros_ip_address" "vlan20" {
  address   = "192.168.2.1/24"
  interface = "VLAN20 - SigenStore"
  network   = "192.168.2.0"
}

resource "routeros_ip_dhcp_client" "wan" {
  interface    = "1-ONT"
  use_peer_dns = false

}

resource "routeros_ipv6_dhcp_client" "wan" {
  interface          = "1-ONT"
  add_default_route  = true
  request            = ["prefix"]
  pool_name          = routeros_ipv6_pool.BRSKPool.name
  pool_prefix_length = 64
  use_interface_duid = true
  use_peer_dns       = false
}

resource "routeros_ipv6_dhcp_server" "server1" {
  address_pool = "static-only"
  name        = "server1"
  interface   = routeros_interface_bridge.bridge.name
  prefix_pool = routeros_ipv6_pool.BRSKPool.name
  disabled    = true
}
