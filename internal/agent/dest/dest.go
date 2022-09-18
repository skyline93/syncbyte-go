package dest

import "github.com/skyline93/syncbyte-go/internal/pkg/types"

type Dest interface {
	Build(destFile string) error
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
