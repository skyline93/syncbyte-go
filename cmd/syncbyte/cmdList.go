package main

import (
	"log"
	"os"

	"github.com/pterm/pterm"
	"github.com/skyline93/syncbyte-go/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "show resource of syncbyte",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

var cmdListBackends = &cobra.Command{
	Use:   "backend",
	Short: "show backend",
	Run: func(cmd *cobra.Command, args []string) {
		PrintS3Backends()
	},
}

var cmdListSource = &cobra.Command{
	Use:   "source",
	Short: "show backup source",
	Run: func(cmd *cobra.Command, args []string) {
		PrintResources()
	},
}

var cmdListBackupJobs = &cobra.Command{
	Use:   "backupjob",
	Short: "show backup job",
	Run: func(cmd *cobra.Command, args []string) {
		PrintBackupJobs()
	},
}

var cmdListBackupSets = &cobra.Command{
	Use:   "backupset",
	Short: "show backup set",
	Run: func(cmd *cobra.Command, args []string) {
		PrintBackupSets()
	},
}

var cmdListRestoreJobs = &cobra.Command{
	Use:   "restorejob",
	Short: "show restore job",
	Run: func(cmd *cobra.Command, args []string) {
		PrintRestoreJobs()
	},
}

var cmdListRestoreResources = &cobra.Command{
	Use:   "restore-resource",
	Short: "show restore resource",
	Run: func(cmd *cobra.Command, args []string) {
		PrintRestoreResources()
	},
}

func init() {
	cmdList.AddCommand(
		cmdListBackends,
		cmdListSource,
		cmdListBackupJobs,
		cmdListBackupSets,
		cmdListRestoreJobs,
		cmdListRestoreResources,
	)

	cmdRoot.AddCommand(cmdList)
}

func PrintS3Backends() {
	v, err := client.Get("backend")
	if err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}

	var rows pterm.TableData
	head := []string{"ID", "EndPoint", "AccessKey", "SecretKey", "Bucket", "StorageType", "DataType"}
	rows = append(rows, head)

	for _, item := range v.([]interface{}) {
		i := S3Backend{}
		if err := Decode(item, &i); err != nil {
			log.Printf("decode body error")
			os.Exit(1)
		}

		row := []string{utils.UintToStr(i.ID), i.EndPoint, i.AccessKey, i.SecretKey, i.Bucket, i.StorageType, string(i.DataType)}
		rows = append(rows, row)
	}

	pterm.DefaultTable.WithHasHeader().WithData(rows).Render()
	pterm.Println()
}

func PrintResources() {
	v, err := client.Get("source")
	if err != nil {
		os.Exit(1)
	}

	var rows pterm.TableData
	head := []string{"ID", "Name", "DBType", "Version", "Server", "Port", "User", "DBName"}
	rows = append(rows, head)

	for _, item := range v.([]interface{}) {
		i := Source{}

		if err := Decode(item, &i); err != nil {
			log.Printf("decode body error")
			os.Exit(1)
		}

		row := []string{utils.UintToStr(i.ID), i.Name, string(i.DbType), i.Version, i.Server, utils.IntToStr(i.Port), i.User, i.DbName}
		rows = append(rows, row)
	}

	pterm.DefaultTable.WithHasHeader().WithData(rows).Render()
	pterm.Println()
}

func PrintBackupJobs() {
	v, err := client.Get("backup/job")
	if err != nil {
		log.Printf("err: %v", err)
		os.Exit(1)
	}

	var rows pterm.TableData
	head := []string{"ID", "StartTime", "EndTime", "Status", "ResourceID", "BackendID"}
	rows = append(rows, head)

	for _, item := range v.([]interface{}) {
		i := BackupJob{}

		if err := Decode(item, &i); err != nil {
			log.Printf("decode body error")
			os.Exit(1)
		}

		row := []string{
			utils.UintToStr(i.ID),
			i.StartTime.Format("2006-01-02 15:04:05"),
			i.EndTime.Format("2006-01-02 15:04:05"),
			string(i.Status),
			utils.UintToStr(i.ResourceID),
			utils.UintToStr(i.BackendID),
		}
		rows = append(rows, row)
	}

	pterm.DefaultTable.WithHasHeader().WithData(rows).Render()
	pterm.Println()
}

func PrintBackupSets() {
	v, err := client.Get("backup/set")
	if err != nil {
		os.Exit(1)
	}

	var rows pterm.TableData
	head := []string{"ID", "DataSetName", "IsValid", "Size", "BackupTime", "IsCompress"}
	rows = append(rows, head)

	for _, item := range v.([]interface{}) {
		i := BackupSet{}

		if err := Decode(item, &i); err != nil {
			log.Printf("decode body error")
			os.Exit(1)
		}

		row := []string{
			utils.UintToStr(i.ID),
			i.DataSetName,
			utils.BoolToStr(i.IsValid),
			utils.IntToStr(i.Size),
			i.BackupTime.Format("2006-01-02 15:04:05"),
			utils.BoolToStr(i.IsCompress),
		}
		rows = append(rows, row)
	}

	pterm.DefaultTable.WithHasHeader().WithData(rows).Render()
	pterm.Println()
}

func PrintRestoreJobs() {
	v, err := client.Get("restore/job")
	if err != nil {
		os.Exit(1)
	}

	var rows pterm.TableData
	head := []string{"ID", "StartTime", "EndTime", "Status", "BackupSetID"}
	rows = append(rows, head)

	for _, item := range v.([]interface{}) {
		i := RestoreJob{}

		if err := Decode(item, &i); err != nil {
			log.Printf("decode body error")
			os.Exit(1)
		}

		row := []string{
			utils.UintToStr(i.ID),
			i.StartTime.Format("2006-01-02 15:04:05"),
			i.EndTime.Format("2006-01-02 15:04:05"),
			string(i.Status),
			utils.UintToStr(i.BackupSetID),
		}
		rows = append(rows, row)
	}

	pterm.DefaultTable.WithHasHeader().WithData(rows).Render()
	pterm.Println()
}

func PrintRestoreResources() {
	v, err := client.Get("restore/resource")
	if err != nil {
		os.Exit(1)
	}

	var rows pterm.TableData
	head := []string{"ID", "Name", "DBType", "Version", "Server", "Port", "User", "DBName", "IsValid", "RestoreTime"}
	rows = append(rows, head)

	for _, item := range v.([]interface{}) {
		i := RestoreDBResource{}

		if err := Decode(item, &i); err != nil {
			log.Printf("decode body error")
			os.Exit(1)
		}

		row := []string{
			utils.UintToStr(i.ID),
			i.Name,
			string(i.DBType),
			i.Version,
			i.Server,
			utils.IntToStr(i.Port),
			i.User,
			i.DBName,
			utils.BoolToStr(i.IsValid),
			i.RestoreTime.Format("2006-01-02 15:04:05"),
		}
		rows = append(rows, row)
	}

	pterm.DefaultTable.WithHasHeader().WithData(rows).Render()
	pterm.Println()
}
