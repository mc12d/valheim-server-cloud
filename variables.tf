variable "yc_cloud" {
  type = string
}

variable "yc_folder" {
  type = string
}

variable "yc_service_account_id" {
  type = string
}

variable "yc_account" {
  type = string
}

variable "valheim_server_supervisor" {
  type = object({
    supervisor_http = bool
    supervisor_http_port = number
    supervisor_http_user = string
    supervisor_http_pass = string
  })
  default = {
    supervisor_http = true
    supervisor_http_user = "super"
    supervisor_http_pass = "repus"
    supervisor_http_port = 8080
  }
}

variable "valheim_server_env" {
  type = object({
    server_name = string
    server_pass = string
    server_public = bool
    world_name = string
    adminlist_ids = list(string)
    tz = string
  })
}

variable "valheim_backup_env" {
    type = object({
        http_port = number
        bucket_id = string
        cron = string
        restore_on_startup = bool
    })
}

variable "vm_access" {
  type = object({
    user = string
    sshPublicKey = string
  })
}

variable "yc_iam_token" {
  type = string
}

variable "s3" {
  type = object({
    bucket = string
    noncurrent_version_expiration_days = number
  })
}
