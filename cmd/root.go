package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	VERSION             = "0.1.0"
	DEFAULT_DAYS_BACK   = 10
	DEFAULT_DATE_FORMAT = "2006-01-02"
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
			fmt.Println("gh-metrics", VERSION)
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
	start := today.AddDate(0, 0, -DEFAULT_DAYS_BACK)

	rootCmd.Flags().StringP("start", "s", start.Format(DEFAULT_DATE_FORMAT), "target start of date range")
	rootCmd.Flags().StringP("end", "e", today.Format(DEFAULT_DATE_FORMAT), "target end of date range")

	rootCmd.Flags().BoolP("csv", "c", false, "print output as CSV")
	rootCmd.Flags().BoolP("version", "v", false, "print current version")
}
