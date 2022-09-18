package options

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/skyline93/syncbyte-go/internal/pkg/types"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var CfgFile string
var Opts *Options

type Options struct {
	ListenAddr string           `json:"listenAddr" mapstructure:"listenAddr"`
	Database   *DatabaseOptions `json:"database" mapstructure:"database"`
}

type DatabaseOptions struct {
	Type     types.DBType `json:"type" mapstructure:"type"`
	Host     string       `json:"host" mapstructure:"host"`
	Port     int          `json:"port" mapstructure:"port"`
	User     string       `json:"user" mapstructure:"user"`
	Password string       `json:"password" mapstructure:"password"`
	DbName   string       `json:"db_name" mapstructure:"dbname"`
	Extra    string       `json:"extra" mapstructure:"extra"`
}

func (opts *DatabaseOptions) Dsn() (dsn string) {
	switch opts.Type {
	case types.MySQL:
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s",
			opts.User, opts.Password, opts.Host, opts.Port, opts.DbName,
		)

		if opts.Extra != "" {
			dsn = fmt.Sprintf("%s?%s", dsn, opts.Extra)
		}
	case types.PostgreSQL:
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d",
			opts.Host, opts.User, opts.Password, opts.DbName, opts.Port,
		)
		if opts.Extra != "" {
			dsn = fmt.Sprintf("%s %s", dsn, opts.Extra)
		}
	case types.SQLite:
		dsn = opts.DbName
	default:
		panic("unsupported database type")
	}

	return dsn
}

func Dialector(opts *DatabaseOptions) gorm.Dialector {
	switch opts.Type {
	case types.MySQL:
		return mysql.Open(opts.Dsn())
	case types.PostgreSQL:
		return postgres.Open(opts.Dsn())
	case types.SQLite:
		return sqlite.Open(opts.Dsn())
	default:
		return nil
	}
}

func InitConfig() {
	if CfgFile != "" {
		viper.SetConfigFile(CfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".syncbyte-engine")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	settings := viper.AllSettings()

	Opts = &Options{}
	mapstructure.Decode(settings, Opts)
}
