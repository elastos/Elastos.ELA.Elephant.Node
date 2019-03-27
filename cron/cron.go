package cron

type Job interface {
	start()
	stop()
}

type CronTable struct {
}
