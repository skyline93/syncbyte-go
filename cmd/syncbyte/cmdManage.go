package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var cmdManage = &cobra.Command{
	Use:   "manage",
	Short: "syncbyte manage of syncbyte",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

var cmdScalePool = &cobra.Command{
	Use:   "scale-pool",
	Short: "worker pool of syncbyte",
	Run: func(cmd *cobra.Command, args []string) {
		if err := scalePool(poolOptions); err != nil {
			log.Printf("err: %v", err)
			os.Exit(1)
		}
	},
}

var cmdListWorker = &cobra.Command{
	Use:   "list-worker",
	Short: "list worker of syncbyte",
	Run: func(cmd *cobra.Command, args []string) {
		PrintPoolWorker()
	},
}

type PoolOptions struct {
	Size int `json:"size"`
}

var poolOptions PoolOptions

func init() {
	cmdManage.AddCommand(cmdScalePool)
	cmdManage.AddCommand(cmdListWorker)
	cmdRoot.AddCommand(cmdManage)

	ef := cmdScalePool.Flags()
	ef.IntVarP(&poolOptions.Size, "size", "s", 0, "worker pool size")
}

func scalePool(poolOpts PoolOptions) error {
	_, err := client.Post("manage/pool", poolOpts)
	if err != nil {
		return err
	}

	return nil
}

func PrintPoolWorker() {
	v, err := client.Get("manage/pool")
	if err != nil {
		log.Printf("err: %v", err)
		os.Exit(1)
	}

	var rows pterm.TableData
	head := []string{"ID", "WorkerID"}
	rows = append(rows, head)

	for n, item := range v.([]interface{}) {
		i := PoolWorkerItem{}

		if err := Decode(item, &i); err != nil {
			log.Printf("decode body error")
			os.Exit(1)
		}

		row := []string{fmt.Sprintf("%d", n+1), i.WorkerID}
		rows = append(rows, row)
	}

	pterm.DefaultTable.WithHasHeader().WithData(rows).Render()
	pterm.Println()
}
