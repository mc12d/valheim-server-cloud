### valheim-server-cloud

This repo provides a way to host your Valheim server using [yandex-cloud](https://cloud.yandex.ru/en) (VM + s3 as backup storage)

It uses [`this docker image`](https://github.com/lloesche/valheim-server-docker) and custom backup agent

Prerequisites:
 - [yandex-cloud CLI](https://cloud.yandex.ru/en/docs/cli/) installed
 - terraform or opentofu installed
 - golang environment (optional)

Steps to launch your server:
 - in yandex cloud UI -> create cloud and folder
 - create and fill `terraform.tfvars` file (use `example.tfvars` for reference)
 - login to `yc` CLI
 - use ```source ./scripts/env.sh``` to export cloud-token to your current  shell
 - `tofu (terraform) plan` -> `tofu apply`
 - bingo! connect to the server via in-game server browser (if `server_public=true` is set), or directly via ip `w.x.y.z:2456` (you can see public ip address in vm information in cloud UI)

Stop the server:
 - start/stop the vm via yandex-cloud UI
 - run `tofu destroy` to delete all cloud resources (including public ip address AND BACKUPS)



Also:
 - navigate to object storage and see your backups. When vm restarts, the most recent backup will be picked up
 - feel free to tweak any settings, you may want to tune additional server settings in `docker-compose.yaml` or vm configuration in `vm.tf`
 - by default, VM is [`preemptible`](https://cloud.yandex.com/en/docs/compute/concepts/preemptible-vm), meaning the cloud may turn it off some time after launch. If you need fully 24/7 server, turn this parameter off in `vm.tf`. _Be careful, it will bump the billing a lot_
 - you can setup discord integration as [described here](https://github.com/lloesche/valheim-server-docker#notify-on-discord), maybe some stuff will be added to this repo in the future

Notes on billing:
 - I will not bring here any numbers, as it varies over time/regions/vm configuration/etc
 - you always can create the resources, see the billing estimates in the UI and destroy everything right after. Use the trial period or just throw in some comfortable amount. Also, use [the calculator](https://cloud.yandex.com/en/prices#calculator)

Valhalla awaits!