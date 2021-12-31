package scheduler

import (
	"app/infrastructure/util/log"

	"github.com/robfig/cron/v3"
)

func StartUp() {
	crontab := cron.New(cron.WithSeconds()) //精确到秒

	log.Info("package cron job start finish")
	log.Info("scheduler start up")
	crontab.Start()
}
