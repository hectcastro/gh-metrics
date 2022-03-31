package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/rickar/cal/v2"
	"github.com/spf13/cobra"
)

const (
	Version           = "0.2.1"
	DefaultDaysBack   = 10
	DefaultDateFormat = "2006-01-02"
)

func WorkdayOnlyWeekdays(day time.Time) bool {
	return day.Weekday() != time.Saturday && day.Weekday() != time.Sunday
}

func WorkdayAllDays(day time.Time) bool {
	return true
}

func WorkdayStart(d time.Time) time.Time {
	year, month, day := d.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, d.Location())
}

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
