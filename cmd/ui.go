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
	DefaultEmptyCell = "--"
	// Default number of search results per query.
	DefaultResultCount = 100
	// Pull request review approved state.
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

// getReadyForReviewOrPrCreatedAt returns when the pull request was
// marked ready for review, or its created date (if it was never in
// a draft state).
func getReadyForReviewOrPrCreatedAt(prCreated string, timelineItems TimelineItems) string {
	if timelineItems.TotalCount == 0 {
		return prCreated
	}

	return timelineItems.Nodes[0].ReadyForReviewEvent.CreatedAt
}

// getTimeToFirstReview returns the time to first review, in hours and
// minutes, for a given PR.
//
//   timeToFirstReview = (readyForReviewAt || prCreatedAt) - firstReviewdAt
//
func (ui *UI) getTimeToFirstReview(author, prCreatedAt string, isDraft bool, timelineItems TimelineItems, reviews Reviews) string {
	// The pull request is still in a draft state, because it has not
	// yet been marked as ready for review.
	if timelineItems.TotalCount == 0 && isDraft {
		return DefaultEmptyCell
	}

	for _, review := range reviews.Nodes {
		if review.Author.Login != author {
			readyForReviewOrPrCreatedAt, _ := time.Parse(time.RFC3339, getReadyForReviewOrPrCreatedAt(prCreatedAt, timelineItems))
			firstReviewedAt, _ := time.Parse(time.RFC3339, review.CreatedAt)

			return formatDuration(ui.subtractTime(firstReviewedAt, readyForReviewOrPrCreatedAt))
		}
	}

	return DefaultEmptyCell
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

// getFirstReviewToLastReview returns the first review to last approving review time, in
// hours and minutes, for a given PR.
//
//   firstReviewToLastReview = lastReviewedAt - firstReviewedAt
//
func (ui *UI) getFirstReviewToLastReview(login string, reviews Reviews) string {
	var nonAuthorReviews ReviewNodes
	for _, review := range reviews.Nodes {
		if review.Author.Login != login {
			nonAuthorReviews = append(nonAuthorReviews, review)
		}
	}

	if len(nonAuthorReviews) == 0 {
		return DefaultEmptyCell
	}

	firstReviewedAt, _ := time.Parse(time.RFC3339, nonAuthorReviews[0].CreatedAt)

	// Iterate in reverse order to get the last approving review
	for i := len(nonAuthorReviews) - 1; i >= 0; i-- {
		if nonAuthorReviews[i].State == ReviewApprovedState {
			lastReviewedAt, _ := time.Parse(time.RFC3339, nonAuthorReviews[i].CreatedAt)
			return formatDuration(ui.subtractTime(lastReviewedAt, firstReviewedAt))
		}
	}

	return DefaultEmptyCell
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
// a set of pull requests determined by the supplied date range, using
// DefaultResultCount.
func (ui *UI) PrintMetrics() string {
	return ui.printMetricsImpl(DefaultResultCount)
}

// printMetricsImpl returns a string representation of the metrics summary
// for a set of pull requests determined by the supplied date range.
func (ui *UI) printMetricsImpl(defaultResultCount int) string {
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

	var gqlQuery MetricsGQLQuery
	var gqlQueryVariables map[string]interface{} = map[string]interface{}{
		"query":       graphql.String(fmt.Sprintf("repo:%s/%s type:pr merged:%s..%s", ui.Owner, ui.Repository, ui.StartDate, ui.EndDate)),
		"resultCount": graphql.Int(defaultResultCount),
		"afterCursor": (*graphql.String)(nil),
	}

	err = client.Query("PullRequests", &gqlQuery, gqlQueryVariables)
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
		"First to Last Review",
		"First Approval to Merge",
	})

	for {
		for _, node := range gqlQuery.Search.Nodes {
			t.AppendRow(table.Row{
				node.PullRequest.Number,
				node.PullRequest.Commits.TotalCount,
				node.PullRequest.Additions,
				node.PullRequest.Deletions,
				node.PullRequest.ChangedFiles,
				ui.getTimeToFirstReview(
					node.PullRequest.Author.Login,
					node.PullRequest.CreatedAt,
					node.PullRequest.IsDraft,
					node.PullRequest.TimelineItems,
					node.PullRequest.Reviews,
				),
				node.PullRequest.Comments.TotalCount,
				node.PullRequest.Participants.TotalCount,
				ui.getFeatureLeadTime(
					node.PullRequest.MergedAt,
					node.PullRequest.Commits,
				),
				ui.getFirstReviewToLastReview(
					node.PullRequest.Author.Login,
					node.PullRequest.Reviews,
				),
				ui.getFirstApprovalToMerge(
					node.PullRequest.Author.Login,
					node.PullRequest.MergedAt,
					node.PullRequest.Reviews,
				),
			})
		}

		if gqlQuery.Search.PageInfo.HasNextPage {
			gqlQueryVariables["afterCursor"] = graphql.String(gqlQuery.Search.PageInfo.EndCursor)
			err = client.Query("PullRequests", &gqlQuery, gqlQueryVariables)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			break
		}
	}

	if ui.CSVFormat {
		return t.RenderCSV()
	}

	return t.Render()
}
