package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var cmdPolicy = &cobra.Command{
	Use:   "policy",
	Short: "backup policy of syncbyte",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

var cmdEnableBackupPolicy = &cobra.Command{
	Use:   "enable",
	Short: "enable backup policy of syncbyte",
	Run: func(cmd *cobra.Command, args []string) {
		if err := enableBackupPolicy(policyOptions.PolicyIDs); err != nil {
			log.Printf("err: %v", err)
			os.Exit(1)
		}
	},
}

var cmdDisableBackupPolicy = &cobra.Command{
	Use:   "disable",
	Short: "disable backup policy of syncbyte",
	Run: func(cmd *cobra.Command, args []string) {
		if err := disableBackupPolicy(policyOptions.PolicyIDs); err != nil {
			log.Printf("err: %v", err)
			os.Exit(1)
		}
	},
}

type PolicyOptions struct {
	PolicyIDs []uint `json:"PolicyIDs"`
}

var policyOptions PolicyOptions

func init() {
	cmdPolicy.AddCommand(cmdEnableBackupPolicy)
	cmdPolicy.AddCommand(cmdDisableBackupPolicy)
	cmdRoot.AddCommand(cmdPolicy)

	ef := cmdEnableBackupPolicy.Flags()
	ef.UintSliceVarP(&policyOptions.PolicyIDs, "backup-policy-id", "i", []uint{}, "The backup policy id")

	df := cmdDisableBackupPolicy.Flags()
	df.UintSliceVarP(&policyOptions.PolicyIDs, "backup-policy-id", "i", []uint{}, "The backup policy id")
}

func enableBackupPolicy(ids []uint) error {
	_, err := client.Post("policy/enable", ids)
	if err != nil {
		return err
	}

	return nil
}

func disableBackupPolicy(ids []uint) error {
	_, err := client.Post("policy/disable", ids)
	if err != nil {
		return err
	}

	return nil
}
