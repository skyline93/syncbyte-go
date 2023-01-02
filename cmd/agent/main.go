package main

import (
	"fmt"
	"os"

	"github.com/skyline93/syncbyte-go/internal/agent"
	"github.com/spf13/cobra"
)

var cmdRoot = &cobra.Command{
	Use:   "syncbyte-agent",
	Short: "syncbyte agent is a backup agent",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

var cmdRun = &cobra.Command{
	Use:   "run",
	Short: "run server of syncbyte-agent",
	Run: func(cmd *cobra.Command, args []string) {
		if err := agent.RunServer(fmt.Sprintf("%s:%d", options.Host, options.Port)); err != nil {
			fmt.Printf("run server failed, err: %v", err)
		}
	},
}

type Options struct {
	Host string
	Port int
}

var options Options

func init() {
	cmdRoot.AddCommand(cmdRun)

	f := cmdRun.PersistentFlags()
	f.StringVarP(&options.Host, "host", "H", "127.0.0.1", "server host")
	f.IntVarP(&options.Port, "port", "p", 50051, "server port")
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
