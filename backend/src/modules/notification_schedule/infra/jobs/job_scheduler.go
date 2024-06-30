package infra_jobs

import (
	"fmt"

	"github.com/robfig/cron"
)

type JobScheduler struct{}

func (jobScheduler JobScheduler) Schedule(jobName string, spec string, job Job) {
	cronJob := cron.New()

	cronJob.AddFunc(spec, func() {
		result := job.Handle()

		if result != nil {
			fmt.Println(result.Error())
		}
	})

	cronJob.Start()
}

type Job interface {
	Handle() error
}
