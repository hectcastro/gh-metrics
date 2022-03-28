package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	gh "github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	DefaultEmptyCell = "--"
)

type UI struct {
	Owner      string
	Repository string
	StartDate  string
	EndDate    string
	CSVFormat  bool
}

func getTimeToFirstReview(prCreatedAtString string, reviews Reviews) string {
	if len(reviews.Nodes) == 0 {
		return DefaultEmptyCell
	}

	firstReviewedAt, _ := time.Parse(time.RFC3339, reviews.Nodes[0].CreatedAt)
	prCreatedAt, _ := time.Parse(time.RFC3339, prCreatedAtString)

	return firstReviewedAt.Sub(prCreatedAt).Round(time.Minute).String()
}

func getFeatureLeadTime(prMergedAtString string, commits Commits) string {
	if len(commits.Nodes) == 0 {
		return DefaultEmptyCell
	}

	prMergedAt, _ := time.Parse(time.RFC3339, prMergedAtString)
	prFirstCommittedAt, _ := time.Parse(time.RFC3339, commits.Nodes[0].Commit.CommittedDate)

	return prMergedAt.Sub(prFirstCommittedAt).Round(time.Minute).String()
}

func getLastReviewToMerge(prMergedAtString string, latestReviews LatestReviews) string {
	if len(latestReviews.Nodes) == 0 {
		return DefaultEmptyCell
	}

	prMergedAt, _ := time.Parse(time.RFC3339, prMergedAtString)
	latestReviewedAt, _ := time.Parse(time.RFC3339, latestReviews.Nodes[0].CreatedAt)

	return prMergedAt.Sub(latestReviewedAt).Round(time.Minute).String()
}

func (ui *UI) PrintMetrics() {
	client, err := gh.GQLClient(
		&api.ClientOptions{
			EnableCache: true,
			CacheTTL:    15 * time.Minute,
			Timeout:     5 * time.Second,
		},
	)
	if err != nil {
		log.Fatal(err)
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
	t.SetOutputMirror(os.Stdout)
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
		"Last Review to Merge",
	})

	for _, node := range gqlQuery.Search.Nodes {
		t.AppendRow(table.Row{
			node.PullRequest.Number,
			node.PullRequest.Commits.TotalCount,
			node.PullRequest.Additions,
			node.PullRequest.Deletions,
			node.PullRequest.ChangedFiles,
			getTimeToFirstReview(
				node.PullRequest.CreatedAt,
				node.PullRequest.Reviews,
			),
			node.PullRequest.Comments.TotalCount,
			node.PullRequest.Participants.TotalCount,
			getFeatureLeadTime(
				node.PullRequest.MergedAt,
				node.PullRequest.Commits,
			),
			getLastReviewToMerge(
				node.PullRequest.MergedAt,
				node.PullRequest.LatestReviews,
			),
		})
	}

	if csvFormat {
		t.RenderCSV()
	} else {
		t.Render()
	}
}
