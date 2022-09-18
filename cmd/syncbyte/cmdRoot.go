package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cmdRoot = &cobra.Command{
	Use:   "syncbyte",
	Short: "syncbyte is a backup tools",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

func init() {
	cobra.OnInitialize(InitConfig, InitClient)
	cmdRoot.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.syncbyte.yaml)")
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
