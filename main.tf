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

module "mikrotik" {
  source         = "./modules/mikrotik"
  ipv6_bridge_ip = var.ipv6_bridge_ip
  ipv6_prefix    = var.ipv6_prefix
  snmp_authentication_password = var.snmp_authentication_password
  snmp_encryption_password      = var.snmp_encryption_password
}

module "graphana" {
  source = "./modules/graphana"
}

variable "tfuserpass" {
  type        = string
  description = "Password for the device"
  sensitive   = true
}

provider "routeros" {
  hosturl  = "https://nev-core-01.local" # env ROS_HOSTURL or MIKROTIK_HOusername    
  username = "tfuser"                    # env ROS_USERNAME or MIKROTIK_USER
  password = var.tfuserpass              # env ROS_PASSWORD or MIKROTIK_PASSWORD
  insecure = true                        # env ROS_INSECURE or MIKROTIK_INSECURE
}




