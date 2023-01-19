// Package cmd implements a command line interface for summarizing
// GitHub pull request metrics.
package cmd

import (
	"fmt"
	"time"

	gh "github.com/cli/go-gh"
	"github.com/rickar/cal/v2"
	"github.com/spf13/cobra"
)

const (
	// Extension version. Displayed when `--version` flag is used.
	Version = "3.0.0"
	// Default number of days in the past to look for pull requests
	// within a repository.
	DefaultDaysBack = 10
	// Default date format to use when displaying dates.
	DefaultDateFormat = "2006-01-02"
)

var (
	// defaultStart is the default start date to query pull requests if none is
	// specified.
	defaultStart string
	// defaultEnd is the default end date to query pull requests if none is
	// specified.
	defaultEnd string
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
	Use:     "gh-metrics",
	Short:   "gh-metrics: provide summary pull request metrics",
	Version: Version,
	RunE: func(cmd *cobra.Command, args []string) error {
		repository, _ := cmd.Flags().GetString("repo")
		startDate, _ := cmd.Flags().GetString("start")
		endDate, _ := cmd.Flags().GetString("end")
		query, _ := cmd.Flags().GetString("query")
		onlyWeekdays, _ := cmd.Flags().GetBool("only-weekdays")
		csvFormat, _ := cmd.Flags().GetBool("csv")

		repo, err := newGHRepo(repository)
		if err != nil {
			return err
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
			Owner:      repo.Owner,
			Repository: repo.Name,
			Host:       repo.Host,
			StartDate:  startDate,
			EndDate:    endDate,
			Query:      query,
			CSVFormat:  csvFormat,
			Calendar:   calendar,
		}

		cmd.Println(ui.PrintMetrics())

		return nil
	},
}

func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	defaultRepo := ""
	currentRepo, _ := gh.CurrentRepository()
	if currentRepo != nil {
		defaultRepo = fmt.Sprintf("%s/%s", currentRepo.Owner(), currentRepo.Name())
	}

	RootCmd.Flags().StringP("repo", "R", defaultRepo, "target repository in '[HOST/]OWNER/REPO' format (defaults to the current working directory's repository)")

	today := time.Now().UTC()
	defaultStart = today.AddDate(0, 0, -DefaultDaysBack).Format(DefaultDateFormat)
	defaultEnd = today.Format(DefaultDateFormat)

	RootCmd.Flags().StringP("start", "s", defaultStart, "target start of date range for merged pull requests")
	RootCmd.Flags().StringP("end", "e", defaultEnd, "target end of date range for merged pull requests")
	RootCmd.Flags().StringP("query", "q", "", "additional query filter for merged pull requests")

	RootCmd.Flags().BoolP("only-weekdays", "w", false, "only include weekdays (M-F) in date range calculations")
	RootCmd.Flags().BoolP("csv", "c", false, "print output as CSV")
}
