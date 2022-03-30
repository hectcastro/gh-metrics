package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	Version           = "0.2.1"
	DefaultDaysBack   = 10
	DefaultDateFormat = "2006-01-02"
)

var RootCmd = &cobra.Command{
	Use:   "gh-metrics",
	Short: "gh-metrics: provide summary pull request metrics",
	Run: func(cmd *cobra.Command, args []string) {
		owner, _ := cmd.Flags().GetString("owner")
		repository, _ := cmd.Flags().GetString("repo")
		startDate, _ := cmd.Flags().GetString("start")
		endDate, _ := cmd.Flags().GetString("end")
		csv, _ := cmd.Flags().GetBool("csv")

		version, _ := cmd.Flags().GetBool("version")
		if version {
			fmt.Println("gh-metrics", Version)
			os.Exit(0)
		}

		ui := &UI{
			Owner:      owner,
			Repository: repository,
			StartDate:  startDate,
			EndDate:    endDate,
			CSVFormat:  csv,
		}

		fmt.Println(ui.PrintMetrics())
	},
}

func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	RootCmd.Flags().StringP("owner", "o", "", "target repository owner")
	RootCmd.MarkFlagRequired("owner")
	RootCmd.Flags().StringP("repo", "r", "", "target repository name")
	RootCmd.MarkFlagRequired("repo")

	today := time.Now().UTC()
	start := today.AddDate(0, 0, -DefaultDaysBack)

	RootCmd.Flags().StringP("start", "s", start.Format(DefaultDateFormat), "target start of date range")
	RootCmd.Flags().StringP("end", "e", today.Format(DefaultDateFormat), "target end of date range")

	RootCmd.Flags().BoolP("csv", "c", false, "print output as CSV")
	RootCmd.Flags().BoolP("version", "v", false, "print current version")
}
