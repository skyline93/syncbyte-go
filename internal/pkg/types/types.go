package types

import (
	"time"
)

type DBType string

const (
	PostgreSQL DBType = "postgresql"
	MySQL      DBType = "mysql"
	SQLite     DBType = "sqlite"
)

type JobStatus string

const (
	Running   JobStatus = "running"
	Successed JobStatus = "successed"
	Failed    JobStatus = "failed"
)

type BackendDataType string

const (
	PGDATA BackendDataType = "pg_data"
)

type DataTypeMapping map[DBType]BackendDataType

var BackendDataTypeMapping = DataTypeMapping{
	PostgreSQL: PGDATA,
}

type SystemReSourceType string

const (
	Backend SystemReSourceType = "backend"
	Source  SystemReSourceType = "source"
	Restore SystemReSourceType = "restore"
)

type BackupScheduleType string

const (
	Cron     BackupScheduleType = "cron"
	Interval BackupScheduleType = "interval"
)

const TimeFormat = "15:04:05"

type LocalTime time.Time

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 2 {
		*t = LocalTime(time.Time{})
		return
	}

	now, err := time.Parse(`"`+TimeFormat+`"`, string(data))
	*t = LocalTime(now)
	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t LocalTime) Value() (string, error) {
	if t.String() == "0001-01-01 00:00:00" {
		return "", nil
	}
	return time.Time(t).Format(TimeFormat), nil
}

func (t *LocalTime) Scan(v interface{}) error {
	tTime, _ := time.Parse(TimeFormat, v.(time.Time).String())
	*t = LocalTime(tTime)
	return nil
}

func (t LocalTime) String() string {
	return time.Time(t).Format(TimeFormat)
}
