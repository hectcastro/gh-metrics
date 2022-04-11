package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/nbio/st"
	"github.com/rickar/cal/v2"
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
                    "author": {
                        "login": "Batman"
                    },
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
                                "author": {
                                    "login": "Joker"
                    },
                                "createdAt": "2022-03-21T15:12:52Z",
                                "state": "APPROVED"
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
		Calendar:   cal.NewBusinessCalendar(),
	}

	st.Assert(t, strings.Contains(ui.PrintMetrics(), "5339"), true)
}

func Test_SearchQuery_WithCSV(t *testing.T) {
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
		CSVFormat:  true,
		Calendar:   cal.NewBusinessCalendar(),
	}

	st.Assert(t, strings.Contains(ui.PrintMetrics(), "5339,4,6,3,1,2m,0,3,1h12m,--,1h9m"), true)
}

func Test_subtractTime_WithinWorkday(t *testing.T) {
	start := time.Date(2022, time.Month(3), 21, 15, 11, 9, 0, time.UTC)
	end := time.Date(2022, time.Month(3), 21, 15, 12, 52, 0, time.UTC)

	ui := &UI{
		Calendar: cal.NewBusinessCalendar(),
	}

	st.Assert(t, ui.subtractTime(end, start).String(), "1m43s")
	st.Assert(t, ui.subtractTime(end, start).String(), "1m43s")
}

func Test_subtractTime_SpanningWeekend(t *testing.T) {
	start := time.Date(2022, time.Month(3), 25, 17, 0, 0, 0, time.UTC)
	end := time.Date(2022, time.Month(3), 28, 0, 0, 0, 0, time.UTC)

	uiWithWeekends := &UI{
		Calendar: &cal.BusinessCalendar{
			WorkdayFunc:      WorkdayAllDays,
			WorkdayStartFunc: WorkdayStart,
			WorkdayEndFunc:   WorkdayEnd,
		},
	}
	st.Assert(t, uiWithWeekends.subtractTime(end, start).String(), "54h59m57s")

	uiWithoutWeekends := &UI{
		Calendar: &cal.BusinessCalendar{
			WorkdayFunc:      WorkdayOnlyWeekdays,
			WorkdayStartFunc: WorkdayStart,
			WorkdayEndFunc:   WorkdayEnd,
		},
	}
	st.Assert(t, uiWithoutWeekends.subtractTime(end, start).String(), "6h59m59s")
}

func Test_formatDuration_LessThanMinute(t *testing.T) {
	st.Assert(t, formatDuration(time.Second*5), DefaultEmptyCell)
}

func Test_formatDuration_MoreThanMinute(t *testing.T) {
	st.Assert(t, formatDuration(time.Minute*5), "5m")
}

func Test_getTimeToFirstReview(t *testing.T) {
	var reviews = Reviews{
		Nodes: ReviewNodes{
			{CreatedAt: "2022-03-22T15:11:09Z"},
		},
	}

	uiWithWeekends := &UI{
		Calendar: &cal.BusinessCalendar{
			WorkdayFunc:      WorkdayAllDays,
			WorkdayStartFunc: WorkdayStart,
			WorkdayEndFunc:   WorkdayEnd,
		},
	}
	st.Assert(t, uiWithWeekends.getTimeToFirstReview("2022-03-21T15:11:09Z", reviews), "24h0m")

	uiWithoutWeekends := &UI{
		Calendar: &cal.BusinessCalendar{
			WorkdayFunc:      WorkdayOnlyWeekdays,
			WorkdayStartFunc: WorkdayStart,
			WorkdayEndFunc:   WorkdayEnd,
		},
	}
	st.Assert(t, uiWithoutWeekends.getTimeToFirstReview("2022-03-21T15:11:09Z", reviews), "24h0m")
}

func Test_getTimeToFirstReview_NoReviews(t *testing.T) {
	var reviews = Reviews{
		Nodes: ReviewNodes{},
	}

	ui := &UI{}
	st.Assert(t, ui.getTimeToFirstReview("2022-03-21T15:11:09Z", reviews), "--")
}

func Test_getFeatureLeadTime(t *testing.T) {
	var commits = Commits{
		TotalCount: 1,
		Nodes: CommitNodes{
			{Commit{CommittedDate: "2022-03-20T15:11:09Z"}},
		},
	}

	uiWithWeekends := &UI{
		Calendar: &cal.BusinessCalendar{
			WorkdayFunc:      WorkdayAllDays,
			WorkdayStartFunc: WorkdayStart,
			WorkdayEndFunc:   WorkdayEnd,
		},
	}
	st.Assert(t, uiWithWeekends.getFeatureLeadTime("2022-03-21T15:11:09Z", commits), "24h0m")

	uiWithoutWeekends := &UI{
		Calendar: &cal.BusinessCalendar{
			WorkdayFunc:      WorkdayOnlyWeekdays,
			WorkdayStartFunc: WorkdayStart,
			WorkdayEndFunc:   WorkdayEnd,
		},
	}
	st.Assert(t, uiWithoutWeekends.getFeatureLeadTime("2022-03-21T15:11:09Z", commits), "15h11m")
}

func Test_getFeatureLeadTime_NoCommits(t *testing.T) {
	var commits = Commits{
		TotalCount: 0,
		Nodes:      CommitNodes{},
	}

	ui := &UI{}
	st.Assert(t, ui.getFeatureLeadTime("2022-03-21T15:11:09Z", commits), "--")
}

func Test_getFirstReviewToLastReview(t *testing.T) {
	var reviews = Reviews{
		Nodes: ReviewNodes{
			{
				Author: Author{
					Login: "Batman",
				},
				CreatedAt: "2022-04-06T15:11:09Z",
				State:     "COMMENTED",
			},
			{
				Author: Author{
					Login: "Joker",
				},
				CreatedAt: "2022-04-06T16:11:09Z",
				State:     "CHANGES_REQUESTED",
			},
			{
				Author: Author{
					Login: "Joker",
				},
				CreatedAt: "2022-04-06T17:11:09Z",
				State:     "APPROVED",
			},
		},
	}

	uiWithWeekends := &UI{
		Calendar: &cal.BusinessCalendar{
			WorkdayFunc:      WorkdayAllDays,
			WorkdayStartFunc: WorkdayStart,
			WorkdayEndFunc:   WorkdayEnd,
		},
	}

	st.Assert(t, uiWithWeekends.getFirstReviewToLastReview("Batman", reviews), "1h0m")
}

func Test_getFirstReviewToLastReview_AuthorReviewLast(t *testing.T) {
	var reviews = Reviews{
		Nodes: ReviewNodes{
			{
				Author: Author{
					Login: "Joker",
				},
				CreatedAt: "2022-04-06T16:11:09Z",
				State:     "CHANGES_REQUESTED",
			},
			{
				Author: Author{
					Login: "Joker",
				},
				CreatedAt: "2022-04-06T17:11:09Z",
				State:     "APPROVED",
			},
			{
				Author: Author{
					Login: "Batman",
				},
				CreatedAt: "2022-04-06T18:11:09Z",
				State:     "COMMENTED",
			},
		},
	}

	uiWithWeekends := &UI{
		Calendar: &cal.BusinessCalendar{
			WorkdayFunc:      WorkdayAllDays,
			WorkdayStartFunc: WorkdayStart,
			WorkdayEndFunc:   WorkdayEnd,
		},
	}

	st.Assert(t, uiWithWeekends.getFirstReviewToLastReview("Batman", reviews), "1h0m")
}

func Test_getFirstReviewToLastReview_ReviewerReviewCommentLast(t *testing.T) {
	var reviews = Reviews{
		Nodes: ReviewNodes{
			{
				Author: Author{
					Login: "Joker",
				},
				CreatedAt: "2022-04-06T16:11:09Z",
				State:     "CHANGES_REQUESTED",
			},
			{
				Author: Author{
					Login: "Joker",
				},
				CreatedAt: "2022-04-06T17:11:09Z",
				State:     "APPROVED",
			},
			{
				Author: Author{
					Login: "Joker",
				},
				CreatedAt: "2022-04-06T18:11:09Z",
				State:     "COMMENTED",
			},
		},
	}

	uiWithWeekends := &UI{
		Calendar: &cal.BusinessCalendar{
			WorkdayFunc:      WorkdayAllDays,
			WorkdayStartFunc: WorkdayStart,
			WorkdayEndFunc:   WorkdayEnd,
		},
	}

	st.Assert(t, uiWithWeekends.getFirstReviewToLastReview("Batman", reviews), "1h0m")
}

func Test_getFirstApprovalToMerge(t *testing.T) {
	var reviews = Reviews{
		Nodes: ReviewNodes{
			{
				Author:    Author{Login: "Batman"},
				CreatedAt: "2022-03-19T15:00:09Z",
				State:     "COMMENTED",
			},
			{
				Author:    Author{Login: "Joker"},
				CreatedAt: "2022-03-20T15:11:09Z",
				State:     "APPROVED",
			},
		},
	}

	uiWithWeekends := &UI{
		Calendar: &cal.BusinessCalendar{
			WorkdayFunc:      WorkdayAllDays,
			WorkdayStartFunc: WorkdayStart,
			WorkdayEndFunc:   WorkdayEnd,
		},
	}
	st.Assert(t, uiWithWeekends.getFirstApprovalToMerge("Batman", "2022-03-21T15:11:09Z", reviews), "24h0m")

	uiWithoutWeekends := &UI{
		Calendar: &cal.BusinessCalendar{
			WorkdayFunc:      WorkdayOnlyWeekdays,
			WorkdayStartFunc: WorkdayStart,
			WorkdayEndFunc:   WorkdayEnd,
		},
	}
	st.Assert(t, uiWithoutWeekends.getFirstApprovalToMerge("Batman", "2022-03-21T15:11:09Z", reviews), "15h11m")
}

func Test_getFirstApprovalToMerge_NoReviews(t *testing.T) {
	var reviews = Reviews{
		Nodes: ReviewNodes{},
	}

	ui := &UI{}
	st.Assert(t, ui.getFirstApprovalToMerge("Batman", "2022-03-21T15:11:09Z", reviews), "--")
}
