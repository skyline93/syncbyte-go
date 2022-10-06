package agent

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/skyline93/gokv"
	"github.com/skyline93/syncbyte-go/internal/pkg/types"
	"github.com/spf13/viper"
)

type Options struct {
	ListenAddr string         `json:"listenAddr" mapstructure:"listenAddr"`
	Local      bool           `json:"local" mapstructure:"local"`
	Backup     BackupOptions  `json:"backup" mapstructure:"backup"`
	Restore    RestoreOptions `json:"restore" mapstructure:"restore"`
}

type BackupOptions struct {
	MountPoint string `json:"mountPoint" mapstructure:"mountPoint"`
}

type RestoreOptions struct {
	MountPoint string `json:"mountPoint" mapstructure:"mountPoint"`
}

var (
	CfgFile string
	Opts    *Options
	Jobs    *gokv.KV = gokv.New()
)

type JobInfo struct {
	Status types.JobStatus
	Size   int64
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
		viper.SetConfigName(".syncbyte-agent")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	settings := viper.AllSettings()

	Opts = &Options{}
	mapstructure.Decode(settings, Opts)
}
