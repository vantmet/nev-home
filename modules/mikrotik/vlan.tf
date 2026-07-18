resource "routeros_interface_vlan" "vlan_10" {
  name       = "VLAN10 - Smart Devices"
  vlan_id    = 10
  interface  = routeros_interface_bridge.bridge.name
}
