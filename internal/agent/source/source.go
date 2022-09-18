package source

import "github.com/skyline93/syncbyte-go/internal/pkg/types"

type Source interface {
	Dump(destFile string) error
}

type Options struct {
	Name     string
	Server   string
	User     string
	Password string
	DBName   string
	Version  string
	DBType   types.DBType
	Port     int
}
