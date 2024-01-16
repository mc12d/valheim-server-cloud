data "yandex_compute_image" "container-optimized-image" {
  family    = "container-optimized-image"
}

resource "yandex_vpc_address" "valheim" {
  folder_id = data.yandex_resourcemanager_folder.valheim.id
  name = "valheim"

  external_ipv4_address {
    zone_id = "ru-central1-a"
  }
}

resource "yandex_vpc_network" "valheim_net" {
  folder_id = data.yandex_resourcemanager_folder.valheim.id
  name = "valheim_net"
}

resource "yandex_vpc_subnet" "subnet" {
  folder_id = data.yandex_resourcemanager_folder.valheim.id
  name = "valheim_subnet"
  network_id     = yandex_vpc_network.valheim_net.id
  v4_cidr_blocks = ["10.0.0.0/24"]
}

resource "yandex_compute_instance" "valheim-vm" {
  folder_id = data.yandex_resourcemanager_folder.valheim.id

  boot_disk {
    initialize_params {
      image_id = data.yandex_compute_image.container-optimized-image.id
      size = 50
      type = "network-hdd"
    }
  }

  network_interface {
    subnet_id = yandex_vpc_subnet.subnet.id
    nat       = true
    nat_ip_address = yandex_vpc_address.valheim.external_ipv4_address[0].address
  }

  scheduling_policy {
    preemptible = true
  }

  allow_stopping_for_update = true

  resources {
    cores  = 2
    core_fraction = 100
    memory = 6
  }

  metadata = {
    docker-compose = templatefile("${path.module}/yaml/docker-compose.yaml", {
      env = var.valheim_server_env,
      backup = var.valheim_backup_env,
      supervisor = var.valheim_server_supervisor,
      s3AccessKeyId = yandex_iam_service_account_static_access_key.sa-static-key.access_key
      s3AccessKeySecret = yandex_iam_service_account_static_access_key.sa-static-key.secret_key
    })
    user-data = templatefile("${path.module}/yaml/cloud-init.yaml", {
      user = var.vm_access.user
      sshPublicKey = file(var.vm_access.sshPublicKey)
    })
    serial-port-enable = 1
  }
}
