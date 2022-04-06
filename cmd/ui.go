package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	gh "github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rickar/cal/v2"
)

const (
	// Default representation of an empty table cell.
	DefaultEmptyCell    = "--"
	ReviewApprovedState = "APPROVED"
)

type UI struct {
	Owner      string
	Repository string
	StartDate  string
	EndDate    string
	CSVFormat  bool
	Calendar   *cal.BusinessCalendar
}

// subtractTime returns the duration t1 - t2, with respect to the
// configured calendar.
func (ui *UI) subtractTime(t1, t2 time.Time) time.Duration {
	return ui.Calendar.WorkHoursInRange(t1, t2)
}

// formatDuration formats a duration in hours and minutes, rounded
// to the nearest minute.
func formatDuration(d time.Duration) string {
	duration := strings.TrimSuffix((d).Round(time.Minute).String(), "0s")

	if len(duration) == 0 {
		return DefaultEmptyCell
	}

	return duration
}

// getTimeToFirstReview returns the time to first review, in hours and
// minutes, for a given PR.
//
//   firstReview = firstReviewAt - prCreatedAt
//
func (ui *UI) getTimeToFirstReview(prCreatedAtString string, reviews Reviews) string {
	if len(reviews.Nodes) == 0 {
		return DefaultEmptyCell
	}

	firstReviewedAt, _ := time.Parse(time.RFC3339, reviews.Nodes[0].CreatedAt)
	prCreatedAt, _ := time.Parse(time.RFC3339, prCreatedAtString)

	return formatDuration(ui.subtractTime(firstReviewedAt, prCreatedAt.UTC()))
}

// getFeatureLeadTime returns the feature lead time, in hours and minutes,
// for a given PR.
//
//   featureLeadTime = prMergedAt - firstCommitAt
//
func (ui *UI) getFeatureLeadTime(prMergedAtString string, commits Commits) string {
	if len(commits.Nodes) == 0 {
		return DefaultEmptyCell
	}

	prMergedAt, _ := time.Parse(time.RFC3339, prMergedAtString)
	prFirstCommittedAt, _ := time.Parse(time.RFC3339, commits.Nodes[0].Commit.CommittedDate)

	return formatDuration(ui.subtractTime(prMergedAt, prFirstCommittedAt))
}

// getFirstApprovalToMerge returns the first approval review to merge time, in
// hours and minutes, for a given PR.
//
//   firstApprovalToMerge = prMergedAt - firstApprovedAt
//
func (ui *UI) getFirstApprovalToMerge(author, prMergedAtString string, reviews Reviews) string {
	for _, review := range reviews.Nodes {
		if review.Author.Login != author && review.State == ReviewApprovedState {
			prMergedAt, _ := time.Parse(time.RFC3339, prMergedAtString)
			firstApprovedAt, _ := time.Parse(time.RFC3339, review.CreatedAt)

			return formatDuration(ui.subtractTime(prMergedAt, firstApprovedAt))
		}
	}

	return DefaultEmptyCell
}

// PrintMetrics returns a string representation of the metrics summary for
// a set of pull requests determined by the supplied date range.
func (ui *UI) PrintMetrics() string {
	client, err := gh.GQLClient(
		&api.ClientOptions{
			EnableCache: true,
			CacheTTL:    15 * time.Minute,
			Timeout:     5 * time.Second,
		},
	)
	if err != nil {
		log.Fatal("To authenticate, please run `gh auth login`.")
	}

	variables := map[string]interface{}{
		"query": graphql.String(fmt.Sprintf("repo:%s/%s type:pr merged:%s..%s", ui.Owner, ui.Repository, ui.StartDate, ui.EndDate)),
	}

	var gqlQuery MetricsGQLQuery
	err = client.Query("PullRequests", &gqlQuery, variables)
	if err != nil {
		log.Fatal(err)
	}

	t := table.NewWriter()
	t.SetStyle(table.StyleLight)

	t.AppendHeader(table.Row{
		"PR",
		"Commits",
		"Additions",
		"Deletions",
		"Changed Files",
		"Time to First Review",
		"Comments",
		"Participants",
		"Feature Lead Time",
		"First Approval to Merge",
	})

	for _, node := range gqlQuery.Search.Nodes {
		t.AppendRow(table.Row{
			node.PullRequest.Number,
			node.PullRequest.Commits.TotalCount,
			node.PullRequest.Additions,
			node.PullRequest.Deletions,
			node.PullRequest.ChangedFiles,
			ui.getTimeToFirstReview(
				node.PullRequest.CreatedAt,
				node.PullRequest.Reviews,
			),
			node.PullRequest.Comments.TotalCount,
			node.PullRequest.Participants.TotalCount,
			ui.getFeatureLeadTime(
				node.PullRequest.MergedAt,
				node.PullRequest.Commits,
			),
			ui.getFirstApprovalToMerge(
				node.PullRequest.Author.Login,
				node.PullRequest.MergedAt,
				node.PullRequest.Reviews,
			),
		})
	}

	if ui.CSVFormat {
		return t.RenderCSV()
	}

	return t.Render()
}
