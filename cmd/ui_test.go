package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

const (
	Owner      = "testOwner"
	Repository = "testRepo"
	StartDate  = "2022-03-18"
	EndDate    = "2022-03-28"

	ResponseJSON = `
{
    "data": {
        "search": {
            "nodes": [
                {
                    "additions": 6,
                    "deletions": 3,
                    "number": 5339,
                    "createdAt": "2022-03-21T15:11:09Z",
                    "changedFiles": 1,
                    "mergedAt": "2022-03-21T16:22:05Z",
                    "participants": {
                        "totalCount": 3
                    },
                    "comments": {
                        "totalCount": 0
                    },
                    "reviews": {
                        "nodes": [
                            {
                                "createdAt": "2022-03-21T15:12:52Z"
                            }
                        ]
                    },
                    "latestReviews": {
                        "nodes": [
                            {
                                "createdAt": "2022-03-21T15:12:52Z"
                            }
                        ]
                    },
                    "commits": {
                        "totalCount": 4,
                        "nodes": [
                            {
                                "commit": {
                                    "committedDate": "2022-03-21T15:09:52Z"
                                }
                            }
                        ]
                    }
                }
            ]
        }
    }
}`
)

type GQLRequest struct {
	Variables struct {
		Query string
	}
}

func gqlSearchQueryMatcher(req *http.Request, ereq *gock.Request) (bool, error) {
	var gqlRequest GQLRequest

	body, err := ioutil.ReadAll(req.Body)
	err = json.Unmarshal(body, &gqlRequest)

	return gqlRequest.Variables.Query == fmt.Sprintf("repo:%s/%s type:pr merged:%s..%s", Owner, Repository, StartDate, EndDate), err
}

func Test_SearchQuery(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com/graphql").
		Post("/").
		MatchType("json").
		AddMatcher(gqlSearchQueryMatcher).
		Reply(200).
		BodyString(ResponseJSON)

	ui := &UI{
		Owner:      Owner,
		Repository: Repository,
		StartDate:  StartDate,
		EndDate:    EndDate,
		CSVFormat:  false,
	}

	st.Assert(t, strings.Contains(ui.PrintMetrics(), "5339"), true)
}

func Test_getTimeToFirstReview(t *testing.T) {
	var reviews = Reviews{
		Nodes: ReviewNodes{
			{CreatedAt: "2022-03-22T15:11:09Z"},
		},
	}

	st.Assert(t, getTimeToFirstReview("2022-03-21T15:11:09Z", reviews), "24h0m0s")
}

func Test_getTimeToFirstReview_NoReviews(t *testing.T) {
	var reviews = Reviews{
		Nodes: ReviewNodes{},
	}

	st.Assert(t, getTimeToFirstReview("2022-03-21T15:11:09Z", reviews), "--")
}

func Test_getFeatureLeadTime(t *testing.T) {
	var commits = Commits{
		TotalCount: 1,
		Nodes: CommitNodes{
			{Commit{CommittedDate: "2022-03-20T15:11:09Z"}},
		},
	}

	st.Assert(t, getFeatureLeadTime("2022-03-21T15:11:09Z", commits), "24h0m0s")
}

func Test_getFeatureLeadTime_NoCommits(t *testing.T) {
	var commits = Commits{
		TotalCount: 0,
		Nodes:      CommitNodes{},
	}

	st.Assert(t, getFeatureLeadTime("2022-03-21T15:11:09Z", commits), "--")
}

func Test_getLastReviewToMerge(t *testing.T) {
	var latestReviews = LatestReviews{
		Nodes: LatestReviewNodes{
			{CreatedAt: "2022-03-20T15:11:09Z"},
		},
	}

	st.Assert(t, getLastReviewToMerge("2022-03-21T15:11:09Z", latestReviews), "24h0m0s")
}

func Test_getLastReviewToMerge_NoReviews(t *testing.T) {
	var latestReviews = LatestReviews{
		Nodes: LatestReviewNodes{},
	}

	st.Assert(t, getLastReviewToMerge("2022-03-21T15:11:09Z", latestReviews), "--")
}
