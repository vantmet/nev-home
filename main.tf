terraform {
  required_providers {
    routeros = {
      source = "terraform-routeros/routeros"
    }
  }
}

variable "tfuserpass" {
  type        = string
  description = "Password for the RDS database instance"
  sensitive   = true
}

provider "routeros" {
  hosturl  = "https://nev-core-01.local" # env ROS_HOSTURL or MIKROTIK_HOusername    
  username = "tfuser"                    # env ROS_USERNAME or MIKROTIK_USER
  password = var.tfuserpass              # env ROS_PASSWORD or MIKROTIK_PASSWORD
  insecure = true                        # env ROS_INSECURE or MIKROTIK_INSECURE
}




