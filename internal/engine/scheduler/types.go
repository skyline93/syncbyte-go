package scheduler

type ScheduleType string

const (
	Cron     ScheduleType = "cron"
	Interval ScheduleType = "interval"
)

type JobType string

const (
	Backup JobType = "backup"
)
