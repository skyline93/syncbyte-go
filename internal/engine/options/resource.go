package options

import "github.com/skyline93/syncbyte-go/internal/pkg/types"


type BackendOption struct {
	EndPoint    string                `yaml:"endpoint"`
	AccessKey   string                `yaml:"accessKey"`
	SecretKey   string                `yaml:"secretKey"`
	Bucket      string                `yaml:"bucket"`
	StorageType string                `yaml:"storageType"`
	DataType    types.BackendDataType `yaml:"dataType"`
}

type SourceOption struct {
	Name     string       `yaml:"name"`
	Server   string       `yaml:"server"`
	Port     int          `yaml:"port"`
	User     string       `yaml:"user"`
	Password string       `yaml:"password"`
	DbName   string       `yaml:"dbname"`
	Extend   string       `yaml:"extend"`
	Version  string       `yaml:"version"`
	DbType   types.DBType `yaml:"type"`
}

type DestOptions struct {
	Name     string
	Server   string
	User     string
	Password string
	DBName   string
	Version  string
	DBType   types.DBType
	Port     int
}

type BackendOptions struct {
	Backends []BackendOption `yaml:"backend"`
}

type SourceOptions struct {
	Sources []SourceOption `yaml:"source"`
}
