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

const DEFAULT_EMPTY_CELL = "--"

type Participants struct {
	TotalCount int
}

type Comments struct {
	TotalCount int
}

type Reviews struct {
	Nodes []struct {
		CreatedAt string
	}
}

type Commits struct {
	TotalCount int
	Nodes      []struct {
		Commit struct {
			CommittedDate string
		}
	}
}

type LatestReviews struct {
	Nodes []struct {
		CreatedAt string
	}
}

type MetricsGQLQuery struct {
	Search struct {
		Nodes []struct {
			PullRequest struct {
				Additions     int
				Deletions     int
				Title         string
				Number        int
				CreatedAt     string
				ChangedFiles  int
				MergedAt      string
				Participants  Participants
				Comments      Comments
				Reviews       Reviews       `graphql:"reviews(first: 1)"`
				LatestReviews LatestReviews `graphql:"latestReviews(first: 1)"`
				Commits       Commits       `graphql:"commits(first: 1)"`
			} `graphql:"... on PullRequest"`
		}
	} `graphql:"search(query: $query, type: ISSUE, last: 50)"`
}

func getTimeToFirstReview(prCreatedAtString string, reviews Reviews) string {
	if len(reviews.Nodes) == 0 {
		return DEFAULT_EMPTY_CELL
	}

	firstReviewedAt, _ := time.Parse(time.RFC3339, reviews.Nodes[0].CreatedAt)
	prCreatedAt, _ := time.Parse(time.RFC3339, prCreatedAtString)

	return firstReviewedAt.Sub(prCreatedAt).Round(time.Minute).String()
}

func getFeatureLeadTime(prMergedAtString string, commits Commits) string {
	if len(commits.Nodes) == 0 {
		return DEFAULT_EMPTY_CELL
	}

	prMergedAt, _ := time.Parse(time.RFC3339, prMergedAtString)
	prFirstCommittedAt, _ := time.Parse(time.RFC3339, commits.Nodes[0].Commit.CommittedDate)

	return prMergedAt.Sub(prFirstCommittedAt).Round(time.Minute).String()
}

func getLastReviewToMerge(prMergedAtString string, latestReviews LatestReviews) string {
	if len(latestReviews.Nodes) == 0 {
		return DEFAULT_EMPTY_CELL
	}

	prMergedAt, _ := time.Parse(time.RFC3339, prMergedAtString)
	latestReviewedAt, _ := time.Parse(time.RFC3339, latestReviews.Nodes[0].CreatedAt)

	return prMergedAt.Sub(latestReviewedAt).Round(time.Minute).String()
}

func printMetrics(owner string, repo string, start string, end string, csvFormat bool) {
	options := api.ClientOptions{
		EnableCache: true,
		CacheTTL:    15 * time.Minute,
		Timeout:     5 * time.Second,
	}

	client, err := gh.GQLClient(&options)
	if err != nil {
		log.Fatal(err)
		return
	}

	variables := map[string]interface{}{
		"query": graphql.String(fmt.Sprintf("repo:%s/%s type:pr merged:%s..%s", owner, repo, start, end)),
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
