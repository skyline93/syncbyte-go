package main

import (
	"fmt"
	"os"

	"github.com/skyline93/syncbyte-go/internal/engine"
	"github.com/spf13/cobra"
)

var cmdRoot = &cobra.Command{
	Use:   "syncbyte-engine",
	Short: "syncbyte engine is a backup engine",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

var cmdBackup = &cobra.Command{
	Use:   "backup",
	Short: "backup",
	Run: func(cmd *cobra.Command, args []string) {
		address := fmt.Sprintf("%s:%d", options.Host, options.Port)

		if err := engine.Backup(address, options.SourcePath, options.MountPoint); err != nil {
			fmt.Printf("backup failed, err: %v", err)
		}
	},
}

type Options struct {
	Host string
	Port int

	SourcePath string
	MountPoint string
}

var options Options

func init() {
	cmdRoot.AddCommand(cmdBackup)

	f := cmdBackup.PersistentFlags()
	f.StringVarP(&options.Host, "host", "H", "127.0.0.1", "server host")
	f.IntVarP(&options.Port, "port", "p", 50051, "server port")
	f.StringVarP(&options.SourcePath, "source", "s", "", "source path")
	f.StringVarP(&options.MountPoint, "dest", "d", "", "dest path")
}

func Execute() {
	if err := cmdRoot.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
