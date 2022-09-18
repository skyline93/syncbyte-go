package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/skyline93/syncbyte-go/internal/pkg/types"
	"github.com/skyline93/syncbyte-go/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var cmdApply = &cobra.Command{
	Use:   "apply",
	Short: "Apply a configuration to a resource by file name",
	Run: func(cmd *cobra.Command, args []string) {
		loadOptionsAndRunFromFile(applyOptions.FileName)
	},
}

type ApplyOptions struct {
	FileName string
}

var (
	applyOptions ApplyOptions
	s3Backend    S3Backend
	source       Source
)

func init() {
	cmdRoot.AddCommand(cmdApply)

	f := cmdApply.PersistentFlags()
	f.StringVarP(&applyOptions.FileName, "file", "f", "", "that contains the configuration to apply")

	cmdApply.MarkFlagRequired("filename")
}

func loadOptionsAndRunFromFile(file string) error {
	opts := SystemResource{}

	if err := utils.DecodeFromFile(file, &opts); err != nil {
		return err
	}

	switch opts.Kind {
	case types.Backend:
		for _, i := range opts.Spec {
			v := S3Backend{}
			if err := mapstructure.Decode(i, &v); err != nil {
				fmt.Printf("err: %v", err)
				continue
			}

			if err := addS3Backend(&s3Backend); err != nil {
				log.Printf("add s3 backend failed, error: %v", err)
				os.Exit(1)
			}
		}
	case types.Source:
		for _, i := range opts.Spec {
			v := Source{}
			if err := mapstructure.Decode(i, &v); err != nil {
				fmt.Printf("err: %v", err)
				continue
			}

			if err := addSource(&source); err != nil {
				log.Printf("add backup source failed, error: %v", err)
				os.Exit(1)
			}
		}
	case types.Restore:
		for _, i := range opts.Spec {
			if err := mapstructure.Decode(i, &restoreOptions); err != nil {
				log.Printf("err: %v", err)
				continue
			}

			if err := startRestore(restoreOptions); err != nil {
				log.Printf("start restore job failed, error: %v", err)
				os.Exit(1)
			}
		}
	}

	return nil
}

func addS3Backend(s3 *S3Backend) error {
	resp, err := client.R().SetBody(s3).Post("backend")
	if err != nil {
		return err
	}

	log.Printf("resp: %v", resp)
	return nil
}

func addSource(source *Source) error {
	resp, err := client.R().SetBody(source).Post("source")
	if err != nil {
		return err
	}

	log.Printf("resp: %v", resp)
	return nil
}
