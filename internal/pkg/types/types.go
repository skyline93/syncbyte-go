package types

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
