package infra_jobs

import (
	"log"

	"github.com/robfig/cron"
)

type JobScheduler struct{}

func (jobScheduler JobScheduler) Schedule(jobName string, spec string, job Job) {
	cronJob := cron.New()

	cronJob.AddFunc(spec, func() {
		err := job.Handle()

		if err != nil {
			log.Fatal(err.Error())
		}
	})

	cronJob.Start()
}

type Job interface {
	Handle() error
}
