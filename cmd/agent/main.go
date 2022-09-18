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
		srv := agent.New()
		if err := srv.Run(); err != nil {
			fmt.Printf("run error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	cobra.OnInitialize(agent.InitConfig)
	cmdRoot.PersistentFlags().StringVarP(&agent.CfgFile, "config", "c", "", "config file (default is $HOME/.syncbyte-agent.yaml)")

	cmdRoot.AddCommand(cmdRun)
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
