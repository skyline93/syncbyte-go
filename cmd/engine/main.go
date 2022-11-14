package main

import (
	"fmt"
	"os"

	"github.com/skyline93/syncbyte-go/internal/engine/options"
	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"github.com/skyline93/syncbyte-go/internal/engine/scheduler"
	"github.com/skyline93/syncbyte-go/internal/engine/webserver"
	"github.com/skyline93/syncbyte-go/internal/pkg/worker"
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

var cmdRun = &cobra.Command{
	Use:   "run",
	Short: "run server of syncbyte-engine",
	Run: func(cmd *cobra.Command, args []string) {
		srv := webserver.New()

		if err := srv.Run(); err != nil {
			fmt.Printf("run error: %v", err)
			os.Exit(1)
		}
	},
}

var cmdRunScheduler = &cobra.Command{
	Use:   "run-scheduler",
	Short: "run scheduler of syncbyte-engine",
	Run: func(cmd *cobra.Command, args []string) {
		worker.InitWorkPool()

		sch := scheduler.New()
		if err := sch.Run(); err != nil {
			panic(err)
		}
	},
}

func init() {
	cobra.OnInitialize(
		options.InitConfig,
		repository.InitDB,
	)
	cmdRoot.AddCommand(cmdRun, cmdRunScheduler)
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
