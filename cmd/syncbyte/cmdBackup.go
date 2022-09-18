package main

import (
	"log"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/skyline93/syncbyte-go/internal/pkg/schema"
	"github.com/spf13/cobra"
)

var cmdBackup = &cobra.Command{
	Use:   "backup",
	Short: "start backup of syncbyte",
	Run: func(cmd *cobra.Command, args []string) {
		if err := startBackup(backupOptions); err != nil {
			log.Printf("start backup failed, error: %v", err)
			os.Exit(1)
		}
	},
}

type BackupOptions struct {
	ResourceID uint `json:"resource_id"`
	IsCompress bool `json:"is_compress"`
}

var backupOptions BackupOptions

func init() {
	cmdRoot.AddCommand(cmdBackup)

	f := cmdBackup.Flags()
	f.UintVarP(&backupOptions.ResourceID, "resource-id", "r", 0, "resource id")
	f.BoolVarP(&backupOptions.IsCompress, "compress", "", true, "compress backup file")
}

func startBackup(backupOpts BackupOptions) error {
	v, err := client.Post("backup", backupOpts)
	if err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}

	bodyData := schema.StartBackupResponse{}
	if err := mapstructure.Decode(v, &bodyData); err != nil {
		log.Printf("decode body error")
		os.Exit(1)
	}

	log.Printf("data: %v", bodyData)
	return nil
}
