resource "docker_image" "telegraf" {
  name = "telegraf:latest"
}

resource "docker_volume" "telegraf_config" {
  name = "telegraf_config"
}

resource "docker_container" "telegraf" {
  name  = "telegraf"
  networks_advanced {
    name = "tig"
  }
  image = docker_image.telegraf.image_id
  restart = "always"

  volumes {
    volume_name      = "${path.cwd}/telegraf"
    container_path = "/etc/telegraf/"
    read_only      = true
  }
  capabilities {
    add  = ["NET_RAW"]
  }
  ports {
    internal = 8125
    external = 8125
    protocol = "udp"
  }
}
