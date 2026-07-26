terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "4.5.0"
    }
  }
}

provider "docker" {
  #host = "unix:///run/podman/podman.sock"
}

resource "docker_network" "tig" {
  name   = "tig"
  ipv6   = true
  driver = "bridge"
  ipam_options = {
    driver = "host-local"
  }
}

# Pulls the image
resource "docker_image" "graphana" {
  name = "grafana/grafana:latest"
}

# want a volume for config
resource "docker_volume" "grafana_config" {
  name = "grafana_config"
}

# Create a container
resource "docker_container" "grapaha-main" {
  image = docker_image.graphana.image_id
  name  = "graphana-main"
  networks_advanced {
    name = docker_network.tig.name
  }
  ports {
    internal = 3000
    external = 3000
  }
  volumes {
    volume_name    = docker_volume.grafana_config.name
    container_path = "/var/lib/graphana"
  }
  restart = "always"
  env     = ["GRAFANA_ADMIN_PASSWORD = /run/secrets/admin"]

}
