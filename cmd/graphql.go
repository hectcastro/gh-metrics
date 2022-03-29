package cmd

type Participants struct {
	TotalCount int
}

type Comments struct {
	TotalCount int
}

type ReviewNodes []struct {
	CreatedAt string
}

type Reviews struct {
	Nodes ReviewNodes
}

type Commit struct {
	CommittedDate string
}

type CommitNodes []struct {
	Commit Commit
}

type Commits struct {
	TotalCount int
	Nodes      CommitNodes
}

type LatestReviewNodes []struct {
	CreatedAt string
}

type LatestReviews struct {
	Nodes LatestReviewNodes
}

type MetricsGQLQuery struct {
	Search struct {
		Nodes []struct {
			PullRequest struct {
				Additions     int
				Deletions     int
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
