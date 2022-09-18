package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/skyline93/syncbyte-go/internal/pkg/schema"
	"github.com/skyline93/syncbyte-go/internal/pkg/types"
	"github.com/spf13/cobra"
)

var cmdRestore = &cobra.Command{
	Use:   "restore",
	Short: "start restore of syncbyte",
	Run: func(cmd *cobra.Command, args []string) {
		if output {
			s := NewSystemResource(types.Restore)
			out, err := s.ToString(restoreOptions)
			if err != nil {
				log.Fatalf("%v", err)
			}

			fmt.Printf("%v", out)
			return
		}

		if err := startRestore(restoreOptions); err != nil {
			log.Printf("start restore failed, error: %v", err)
			os.Exit(1)
		}
	},
}

type RestoreOptions struct {
	BackupSetID uint   `json:"backup_set_id" yaml:"backupSetID"`
	Name        string `json:"name" yaml:"name"`
	Server      string `json:"server" yaml:"server"`
	Port        int    `json:"port" yaml:"port"`
	User        string `json:"user" yaml:"user"`
	Password    string `json:"password" yaml:"password"`
	DbName      string `json:"dbname" yaml:"dbName"`
	Version     string `json:"version" yaml:"version"`
	DbType      string `json:"db_type" yaml:"dbType"`
}

var output bool

var restoreOptions *RestoreOptions

func init() {
	cmdRoot.AddCommand(cmdRestore)

	restoreOptions = &RestoreOptions{}

	f := cmdRestore.Flags()
	f.UintVarP(&restoreOptions.BackupSetID, "backupset-id", "", 1, "backup set id of restore")
	f.StringVarP(&restoreOptions.DbType, "type", "t", "postgresql", "source type")
	f.StringVarP(&restoreOptions.Name, "name", "n", "source", "source name")
	f.StringVarP(&restoreOptions.Server, "server", "s", "127.0.0.1", "database host, ip or domain name")
	f.IntVarP(&restoreOptions.Port, "port", "p", 5432, "database port")
	f.StringVarP(&restoreOptions.User, "user", "u", "syncbyte", "database user name")
	f.StringVarP(&restoreOptions.Password, "password", "", "passwordSYNCBYTE", "database password")
	f.StringVarP(&restoreOptions.DbName, "dbname", "", "syncbytedb", "database name")
	f.StringVarP(&restoreOptions.Version, "version", "v", "14.5", "database version")
	f.BoolVarP(&output, "output", "o", true, "output to stdout")
}

func startRestore(restoreOpts *RestoreOptions) error {
	v, err := client.Post("restore", restoreOpts)
	if err != nil {
		log.Printf("err: %v", err)
		os.Exit(1)
	}

	bodyData := schema.StartRestoreResponse{}
	if err := mapstructure.Decode(v, &bodyData); err != nil {
		log.Printf("decode body error")
		os.Exit(1)
	}

	log.Printf("data: %v", bodyData)
	return nil
}
