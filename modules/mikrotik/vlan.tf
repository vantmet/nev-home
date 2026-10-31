resource "routeros_interface_vlan" "vlan20" {
  name      = "VLAN20 - SigenStore"
  vlan_id   = 20
  interface = routeros_interface_bridge.bridge.name
}

resource "routeros_interface_bridge_vlan" "vlan20" {
  bridge   = routeros_interface_bridge.bridge.name
  vlan_ids = ["20"]
  untagged = ["7-P24", "8-P23"]
}
