# yandex cloud
yc_account = "yandex-account-id"
yc_service_account_id = "desired-service-account-id"
yc_cloud = "valheim-cloud"
yc_folder = "valheim"

# this user will be created in vm
vm_access = {
    user = "user"
    sshPublicKey = "~/.ssh/id_rsa.pub"
}

s3 = {
    bucket = "valheim-backup"
    noncurrent_version_expiration_days = 7
}

# these vars define env for docker images
valheim_server_supervisor = {
    supervisor_http = true
    supervisor_http_port = 8080
    supervisor_http_user = "username"
    supervisor_http_pass = "password"
}

valheim_server_env = {
    server_name = "name"
    server_pass = "pass"
    world_name = "world"
    server_public = true
    adminlist_ids = ["steamId1", "steamId2"]
    tz = "RU"
}

# Available backup-agent httpserver endpoints:
#  - POST /backup/make --> uploads current worlds_local to s3 (zip-compressed)
#  - POST /backup/restore(?version=ABCDE) --> replaces worlds_local contents with downloaded backup,
#     newest version if not specified
#
# Note: valheim-server itself updates worlds_local content every 20 mins (or via manual save or sleep)
#   10 mins chosen to half maximum backup delay
valheim_backup_env = {
    http_port = 9080
    bucket_id = "valheim-backup"
    cron = "*/10 * * * *"
    # replace worlds_local content with latest backup on startup
    restore_on_startup = true
}

