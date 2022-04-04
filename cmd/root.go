// Package cmd implements a command line interface for summarizing
// GitHub pull request metrics.
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/rickar/cal/v2"
	"github.com/spf13/cobra"
)

const (
	// Extension version. Displayed when `--version` flag is used.
	Version = "0.3.1"
	// Default number of days in the past to look for pull requests
	// within a repository.
	DefaultDaysBack = 10
	// Default date format to use when displaying dates.
	DefaultDateFormat = "2006-01-02"
)

// WorkdayOnlyWeekdays returns true if the given day is a weekday,
// otherwise returns false.
func WorkdayOnlyWeekdays(d time.Time) bool {
	return d.Weekday() != time.Saturday && d.Weekday() != time.Sunday
}

// WorkdayAllDays returns true regardless of the given day, as currently
// all days are considered workdays.
func WorkdayAllDays(d time.Time) bool {
	return true
}

// WorkdayStart determines the beginning of a workday by returning the
// same day, but at the first second.
func WorkdayStart(d time.Time) time.Time {
	year, month, day := d.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, d.Location())
}

// WorkdayEnd determines the end of a workday by returning the same day,
// but at the last second.
func WorkdayEnd(d time.Time) time.Time {
	year, month, day := d.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, d.Location())
}

var RootCmd = &cobra.Command{
	Use:   "gh-metrics",
	Short: "gh-metrics: provide summary pull request metrics",
	Run: func(cmd *cobra.Command, args []string) {
		owner, _ := cmd.Flags().GetString("owner")
		repository, _ := cmd.Flags().GetString("repo")
		startDate, _ := cmd.Flags().GetString("start")
		endDate, _ := cmd.Flags().GetString("end")
		onlyWeekdays, _ := cmd.Flags().GetBool("only-weekdays")
		csvFormat, _ := cmd.Flags().GetBool("csv")

		version, _ := cmd.Flags().GetBool("version")
		if version {
			fmt.Println("gh-metrics", Version)
			os.Exit(0)
		}

		var workdayFunc cal.WorkdayFn
		if onlyWeekdays {
			workdayFunc = WorkdayOnlyWeekdays
		} else {
			workdayFunc = WorkdayAllDays
		}

		calendar := &cal.BusinessCalendar{
			WorkdayFunc:      workdayFunc,
			WorkdayStartFunc: WorkdayStart,
			WorkdayEndFunc:   WorkdayEnd,
		}

		ui := &UI{
			Owner:      owner,
			Repository: repository,
			StartDate:  startDate,
			EndDate:    endDate,
			CSVFormat:  csvFormat,
			Calendar:   calendar,
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

	RootCmd.Flags().BoolP("only-weekdays", "w", false, "only include weekdays (M-F) in date range calculations")
	RootCmd.Flags().BoolP("csv", "c", false, "print output as CSV")
	RootCmd.Flags().BoolP("version", "v", false, "print current version")
}
