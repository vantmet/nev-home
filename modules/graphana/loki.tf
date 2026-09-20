resource "docker_image" "loki" {
  name = "grafana/loki:latest"
}
resource "docker_image" "alloy" {
  name = "grafana/alloy:latest"
}

resource "docker_image" "rsyslog" {
  name = "rsyslog/rsyslog:latest"
}

resource "docker_container" "loki" {
  name = "loki"
  command = [
    "-config.file=/mnt/config/loki-config.yaml",
  ]
  image = docker_image.loki.image_id
  ports {
    internal = 3100
    external = 3100
  }
  volumes {
    volume_name    = "${path.cwd}/loki/"
    container_path = "/mnt/config"
    read_only      = true
  }
  networks_advanced {
    name = "tig"
  }
}

resource "docker_container" "alloy" {
  name = "alloy"
  command = [
    "run",
    "--server.http.listen-addr=0.0.0.0:12345",
    "--storage.path=/var/lib/alloy/data",
    "/etc/alloy/config.alloy"
  ]
  image = docker_image.alloy.image_id
  ports {
    internal = 12347
    external = 12347
    protocol = "tcp"
  }
  ports {
    internal = 12346
    external = 12346
    protocol = "udp"
  }
  volumes {
    volume_name    = "${path.cwd}/alloy/"
    container_path = "/etc/alloy"
    read_only      = true
  }
  volumes {
    volume_name = "/run/user/1000/podman/podman.sock"
    container_path = "/run/user/1000/podman/podman.sock"
  }
  networks_advanced {
    name = "tig"
  }
}

resource "docker_container" "rsyslog" {
  name = "rsyslog"
  image = docker_image.rsyslog.image_id
  ports {
    internal = 10514
    external = 10514
    protocol = "tcp"
  }
  ports {
    internal = 10514
    external = 10514
    protocol = "udp"
  }
  volumes {
    volume_name    = "${path.cwd}/rsyslog/rsyslog.conf"
    container_path = "/etc/rsyslog.conf"
    read_only      = true
  }
  networks_advanced {
    name = "tig"
  }
}
