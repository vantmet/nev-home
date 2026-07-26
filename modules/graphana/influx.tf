# Pull the image
resource "docker_image" "influxdb" {
  name = "influxdb:latest"
}

# data and config volumes
resource "docker_volume" "influxdb_data" {
  name = "influxdb_data"
}
resource "docker_volume" "influxdb_config" {
  name = "influxdb_config"
}

# container
resource "docker_container" "influxdb" {
  name = "influxdb"
  networks_advanced {
    name = docker_network.tig.name
  }
  image = docker_image.influxdb.image_id

  # Map the data and config volumes to the container
  volumes {
    container_path = "/var/lib/influxdb2"
    volume_name    = docker_volume.influxdb_data.name
  }
  volumes {
    container_path = "/etc/influxdb2"
    volume_name    = docker_volume.influxdb_config.name
  }

  # Expose the necessary ports
  ports {
    internal = 8086
    external = 8086
  }
  env = [
    "DOCKER_INFLUXDB_INIT_MODE = setup",
    "DOCKER_INFLUXDB_INIT_USERNAME = admin",
    "DOCKER_INFLUXDB_INIT_PASSWORD = adminpassword123",
    "DOCKER_INFLUXDB_INIT_ORG = myorg",
    "DOCKER_INFLUXDB_INIT_BUCKET = telegraf",
    "DOCKER_INFLUXDB_INIT_ADMIN_TOKEN = my-super-secret-token",
    "DOCKER_INFLUXDB_INIT_RETENTION = 30d"
  ]
}
