resource "yandex_iam_service_account_static_access_key" "sa-static-key" {
  service_account_id = yandex_iam_service_account.sa.id
  description        = "static access key for object storage"
}

resource "yandex_storage_bucket" "valheim-backup" {
  folder_id = data.yandex_resourcemanager_folder.valheim.id

  access_key = yandex_iam_service_account_static_access_key.sa-static-key.access_key
  secret_key = yandex_iam_service_account_static_access_key.sa-static-key.secret_key

  bucket = var.s3.bucket
  acl = "private"

  versioning {
    enabled = true
  }
  
  lifecycle_rule {
    id = "expiration"
    enabled = true
    
    noncurrent_version_expiration {
      days = var.s3.noncurrent_version_expiration_days
    }
  }
}
