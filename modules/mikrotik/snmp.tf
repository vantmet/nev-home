variable "snmp_authentication_password" {
  description = "SNMP authentication password"
  type        = string
  sensitive   = true
}
variable "snmp_encryption_password" {
  description = "SNMP encryption password"
  type        = string
  sensitive   = true
}

resource "routeros_snmp" "snmp" {
  enabled = true
}

resource "routeros_snmp_community" "public" {
  name      = "public"
  disabled = true
}
resource "routeros_snmp_community" "grafana" {
  name = "Grafana"
  authentication_password = var.snmp_authentication_password
  authentication_protocol = "MD5"
  disabled                = false
  encryption_password     = var.snmp_encryption_password
  encryption_protocol     = "DES"
  read_access             = true
  write_access            = false
}
