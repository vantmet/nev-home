resource "routeros_interface_bridge" "bridge" {
  comment        = "defconf"
  name           = "bridge"
  vlan_filtering = false
}

resource "routeros_interface_bridge_port" "bridge_ports" {
  for_each = {
    "2-BRSK"       = { comment = "", pvid = "1" }
    "3-Powerline"  = { comment = "", pvid = "1" }
    "4-Desk"       = { comment = "", pvid = "1" }
    "5-Hive"       = { comment = "", pvid = "1" }
    "6-New RPi"    = { comment = "", pvid = "1" }
    "7-E Room"     = { comment = "", pvid = "1" }
    "8-OldRPi"     = { comment = "", pvid = "1" }
    "sfp-sfpplus1" = { comment = "", pvid = "1" }
  }
  bridge    = routeros_interface_bridge.bridge.name
  interface = each.key
  comment   = each.value.comment
  pvid      = each.value.pvid
}
