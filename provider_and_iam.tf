terraform {
  required_providers {
    yandex = {
      source = "yandex-cloud/yandex"
    }
  }
  required_version = ">= 0.13"
}

provider "yandex" {
  token = var.yc_iam_token
  zone = "ru-central1-a"
  cloud_id = var.yc_cloud
}

data "yandex_iam_user" "me" {
  login = var.yc_account
}

data "yandex_resourcemanager_cloud" "valheim_cloud" {
  name = var.yc_cloud
}

data "yandex_resourcemanager_folder" "valheim" {
  name = var.yc_folder
  cloud_id = data.yandex_resourcemanager_cloud.valheim_cloud.id
}

resource "yandex_iam_service_account" "sa" {
  folder_id = data.yandex_resourcemanager_folder.valheim.id
  name      = var.yc_service_account_id
}

resource "yandex_resourcemanager_cloud_iam_binding" "cloud_binding_sa" {
  members   = [
    "userAccount:${data.yandex_iam_user.me.id}",
    "serviceAccount:${yandex_iam_service_account.sa.id}"
  ]
  role      = "admin"
  cloud_id  = data.yandex_resourcemanager_cloud.valheim_cloud.id
}

resource "yandex_resourcemanager_folder_iam_binding" "folder_binding" {
  folder_id = data.yandex_resourcemanager_folder.valheim.id
  members   = [
    "userAccount:${data.yandex_iam_user.me.id}",
    "serviceAccount:${yandex_iam_service_account.sa.id}"
  ]
  role      = "admin"
}
