package cmd

type Participants struct {
	TotalCount int
}

type Comments struct {
	TotalCount int
}

type Author struct {
	Login string
}

type ReviewNodes []struct {
	Author    Author
	CreatedAt string
	State     string
}

type Reviews struct {
	TotalCount int
	Nodes      ReviewNodes
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

type ReadyForReviewEvent struct {
	CreatedAt string
}

type TimelineItemNodes []struct {
	ReadyForReviewEvent ReadyForReviewEvent `graphql:"... on ReadyForReviewEvent"`
}

type TimelineItems struct {
	TotalCount int
	Nodes      TimelineItemNodes
}

type MetricsGQLQuery struct {
	Search struct {
		Nodes []struct {
			PullRequest struct {
				Author        Author
				Additions     int
				Deletions     int
				Number        int
				CreatedAt     string
				ChangedFiles  int
				IsDraft       bool
				MergedAt      string
				Participants  Participants
				Comments      Comments
				Reviews       Reviews       `graphql:"reviews(first: 50, states: [APPROVED, CHANGES_REQUESTED, COMMENTED])"`
				Commits       Commits       `graphql:"commits(first: 1)"`
				TimelineItems TimelineItems `graphql:"timelineItems(first: 1, itemTypes: [READY_FOR_REVIEW_EVENT])"`
			} `graphql:"... on PullRequest"`
		}
	} `graphql:"search(query: $query, type: ISSUE, last: 50)"`
}
