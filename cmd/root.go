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

var rootCmd = &cobra.Command{
	Use:   "gh-metrics",
	Short: "gh-metrics: provide summary pull request metrics",
	Run: func(cmd *cobra.Command, args []string) {
		owner, _ := cmd.Flags().GetString("owner")
		repo, _ := cmd.Flags().GetString("repo")
		start, _ := cmd.Flags().GetString("start")
		end, _ := cmd.Flags().GetString("end")
		csv, _ := cmd.Flags().GetBool("csv")

		version, _ := cmd.Flags().GetBool("version")
		if version {
			fmt.Println("gh-metrics", Version)
			os.Exit(0)
		}

		printMetrics(owner, repo, start, end, csv)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().StringP("owner", "o", "", "target repository owner")
	rootCmd.MarkFlagRequired("owner")
	rootCmd.Flags().StringP("repo", "r", "", "target repository name")
	rootCmd.MarkFlagRequired("repo")

	today := time.Now().UTC()
	start := today.AddDate(0, 0, -DefaultDaysBack)

	rootCmd.Flags().StringP("start", "s", start.Format(DefaultDateFormat), "target start of date range")
	rootCmd.Flags().StringP("end", "e", today.Format(DefaultDateFormat), "target end of date range")

	rootCmd.Flags().BoolP("csv", "c", false, "print output as CSV")
	rootCmd.Flags().BoolP("version", "v", false, "print current version")
}
