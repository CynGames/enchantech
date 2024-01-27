package di

import (
	"enchantech-codex/src/feeds/service"
	"enchantech-codex/src/utils"
	"github.com/go-co-op/gocron/v2"
	"time"
)

type CronJob func(*service.FeedService) error

func InitializeScheduler(feedService *service.FeedService) {
	scheduler, err := gocron.NewScheduler()

	utils.ErrorPanicPrinter(err, true)
	AddJob(scheduler, time.Hour, CronJobLogic, feedService)

	//scheduler.Start()
}

func AddJob(scheduler gocron.Scheduler, duration time.Duration, job CronJob, feedService *service.FeedService) {
	println("Adding job to scheduler...")

	_, err := scheduler.NewJob(
		gocron.DurationJob(duration),
		gocron.NewTask(func() error { return job(feedService) }),
	)

	utils.ErrorPanicPrinter(err, true)
}

func CronJobLogic(feedService *service.FeedService) error {
	println("CRON article updater running...")

	return feedService.GetRSSXMLContent()
}
