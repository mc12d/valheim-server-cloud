package app

import (
	"github.com/robfig/cron/v3"
)

func BackupCron(logger cron.Logger) *cron.Cron {
	return cron.New(cron.WithLogger(logger))

}
