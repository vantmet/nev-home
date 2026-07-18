import {
  to = module.mikrotik.routeros_interface_bridge.bridge
  id = "*B"
}
import {
  for_each = {
    "2-BRSK"       = { id = "*0" }
    "3-Powerline"  = { id = "*1" }
    "4-Desk"       = { id = "*2" }
    "5-Hive"       = { id = "*3" }
    "6-New RPi"    = { id = "*4" }
    "7-E Room"     = { id = "*5" }
    "8-OldRPi"     = { id = "*6" }
    "sfp-sfpplus1" = { id = "*7" }
  }
  to = module.mikrotik.routeros_interface_bridge_port.bridge_ports[each.key]
  id = each.value.id
}
import {
  to = module.mikrotik.routeros_ip_pool.dhcp
  id = "*1"
}
import {
  to = module.mikrotik.routeros_ip_dhcp_server_network.dhcp
  id = "*1"
}
import {
  to = module.mikrotik.routeros_ip_dhcp_server.defconf
  id = "*1"
}
import {
  for_each = {
    "NRPI"    = { id = "*C" }
    "BRSK"    = { id = "*115" }
    "PLUTO"   = { id = "*1B8" }
    "VENUS"   = { id = "*FCB" }
    "PRINTER" = { id = "*14C9" }
    "MMESH"   = { id = "*1739" }
    "OMESH"   = { id = "*174C" }
    "JMESH"   = { id = "*1823" }
  }
  to = module.mikrotik.routeros_ip_dhcp_server_lease.static[each.key]
  id = each.value.id
}
import {
  to = module.mikrotik.routeros_ip_vrf.vrf_main
  id = "*0"
}
import {
  to = module.mikrotik.routeros_ipv6_pool.BRSKPool
  id = "*1"
}
import {
  to = module.mikrotik.routeros_ipv6_address.lanv6
  id = "*C"
}
import {
  to = module.mikrotik.routeros_ip_address.lan
  id = "*1"
}
import {
  to = module.mikrotik.routeros_ip_dhcp_client.wan
  id = "*1"
}
import {
  to = module.mikrotik.routeros_system_certificate.local-root-ca-cert
  id = "*2"
}
