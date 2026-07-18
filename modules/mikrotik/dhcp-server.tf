resource "routeros_ip_pool" "dhcp" {
  name   = "default-dhcp"
  ranges = ["192.168.1.10-192.168.1.254"]
}

resource "routeros_ip_dhcp_server_network" "dhcp" {
  address    = "192.168.1.0/24"
  netmask    = "24"
  gateway    = "192.168.1.1"
  dns_server = ["192.168.1.21"]
  ntp_server = ["192.168.1.1"]
  domain     = "anferny.local"
}

resource "routeros_ip_dhcp_server" "defconf" {
  name                      = "defconf"
  address_pool              = routeros_ip_pool.dhcp.name
  interface                 = routeros_interface_bridge.bridge.name
  add_arp                   = true
  dynamic_lease_identifiers = "client-mac,client-id"
}

resource "routeros_ip_dhcp_server_lease" "static" {
  for_each = {
    "NRPI"    = { mac = "D8:3A:DD:E7:EE:54", addr = "192.168.1.21", cid = "1:d8:3a:dd:e7:ee:54" }
    "VENUS"   = { mac = "38:7C:76:D2:AC:54", addr = "192.168.1.58", cid = "ff:76:d2:ac:54:0:1:0:1:30:ed:34:e3:38:7c:76:d2:ac:54" }
    "PRINTER" = { mac = "30:05:5C:A7:45:CF", addr = "192.168.1.82", cid = null }
    "DON-M"   = { mac = "28:05:A5:88:77:C7", addr = "192.168.1.249", cid = null }
    "PLUTO"   = { mac = "E4:54:E8:33:3A:55", addr = "192.168.1.222", cid = null }
    "BRSK"    = { mac = "D4:92:5E:79:36:FB", addr = "192.168.1.250", cid = null }
    "MMESH"   = { mac = "D4:92:5E:70:AF:91", addr = "192.168.1.251", cid = null }
    "OMESH"   = { mac = "08:C7:F5:89:B1:4D", addr = "192.168.1.252", cid = null }
    "JMESH"   = { mac = "08:C7:F5:89:C3:D1", addr = "192.168.1.253", cid = null }
  }

  mac_address      = each.value.mac
  address          = each.value.addr
  client_id        = each.value.cid
  server           = routeros_ip_dhcp_server.defconf.name
  comment          = each.key
  always_broadcast = true
}
