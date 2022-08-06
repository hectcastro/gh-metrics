package cmd

import (
	"encoding/json"
	"fmt"
	"io"
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
            "pageInfo": {
                "hasNextPage": false,
                "endCursor": "Y3Vyc29yOjI="
            },
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
                    "isDraft": false,
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
                                "state": "COMMENTED"
                            },
                            {
                                "author": {
                                    "login": "Joker"
                                },
                                "createdAt": "2022-03-22T15:12:52Z",
                                "state": "APPROVED"
                            }
                        ]
                    },
                    "commits": {
                        "totalCount": 1,
                        "nodes": [
                            {
                                "commit": {
                                    "committedDate": "2022-03-21T15:09:52Z"
                                }
                            }
                        ]
                    },
                    "timelineItems": {
                        "totalCount": 1,
                        "nodes": [
                            {
                                "createdAt": "2022-03-15T03:46:20Z"
                            }
                        ]
                    }
                },
                {
                    "author": {
                        "login": "Batman"
                    },
                    "additions": 12,
                    "deletions": 6,
                    "number": 5340,
                    "createdAt": "2022-03-22T15:11:09Z",
                    "changedFiles": 2,
                    "isDraft": false,
                    "mergedAt": "2022-03-22T16:22:05Z",
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
                                "createdAt": "2022-03-22T15:12:52Z",
                                "state": "COMMENTED"
                            },
                            {
                                "author": {
                                    "login": "Joker"
                                },
                                "createdAt": "2022-03-23T15:12:52Z",
                                "state": "APPROVED"
                            }
                        ]
                    },
                    "commits": {
                        "totalCount": 1,
                        "nodes": [
                            {
                                "commit": {
                                    "committedDate": "2022-03-22T15:09:52Z"
                                }
                            }
                        ]
                    },
                    "timelineItems": {
                        "totalCount": 1,
                        "nodes": [
                            {
                                "createdAt": "2022-03-16T03:46:20Z"
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

	var body, err = io.ReadAll(req.Body)
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

	have := ui.PrintMetrics()

	st.Assert(t, strings.Contains(have, "5339"), true)
	st.Assert(t, strings.Contains(have, "5340"), true)
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

	have := ui.PrintMetrics()

	st.Assert(t, strings.Contains(have, "5339,1,6,3,1,38:13,0,3,01:12,08:00,06:51"), true)
	st.Assert(t, strings.Contains(have, "5340,1,12,6,2,38:13,0,3,01:12,08:00,06:51"), true)
}

func Test_SearchQuery_WithPagination(t *testing.T) {
	defer gock.Off()

	responseJSONWithPagination := strings.ReplaceAll(
		strings.Clone(ResponseJSON),
		"\"hasNextPage\": false,",
		"\"hasNextPage\": true,",
	)

	gock.New("https://api.github.com/graphql").
		Post("/").
		MatchType("json").
		AddMatcher(gqlSearchQueryMatcher).
		Reply(200).
		BodyString(responseJSONWithPagination)

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

	have := ui.printMetricsImpl(1)

	st.Assert(t, strings.Contains(have, "5339,1,6,3,1,38:13,0,3,01:12,08:00,06:51"), true)
	st.Assert(t, strings.Contains(have, "5340,1,12,6,2,38:13,0,3,01:12,08:00,06:51"), true)
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
	st.Assert(t, formatDuration(time.Second*5, false), DefaultEmptyCell)
}

func Test_formatDuration_LessThanMinuteWithCSV(t *testing.T) {
	st.Assert(t, formatDuration(time.Second*5, true), "00:00")
}

func Test_formatDuration_MoreThanMinute(t *testing.T) {
	st.Assert(t, formatDuration(time.Minute*5, false), "5m")
}

func Test_formatDuration_MoreThanMinuteWithCSV(t *testing.T) {
	st.Assert(t, formatDuration(time.Minute*5, true), "00:05")
}

func Test_getReadyForReviewOrPrCreatedAt_prCreatedAt(t *testing.T) {
	st.Assert(t, getReadyForReviewOrPrCreatedAt("2022-03-21T15:11:09Z", TimelineItems{
		TotalCount: 0,
	}) == "2022-03-21T15:11:09Z", true)
}

func Test_getReadyForReviewOrPrCreatedAt_readyForReviewAt(t *testing.T) {
	st.Assert(t, getReadyForReviewOrPrCreatedAt("2022-03-21T15:11:09Z", TimelineItems{
		TotalCount: 1,
		Nodes: TimelineItemNodes{
			{ReadyForReviewEvent{CreatedAt: "2022-03-22T15:11:09Z"}},
		}}) == "2022-03-22T15:11:09Z", true)
}

func Test_getTimeToFirstReview(t *testing.T) {
	var timelineItems = TimelineItems{
		TotalCount: 1,
		Nodes: TimelineItemNodes{
			{ReadyForReviewEvent{CreatedAt: "2022-03-21T15:11:09Z"}},
		},
	}
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
	st.Assert(t, uiWithWeekends.getTimeToFirstReview("Batman", "", false, timelineItems, reviews), "24h0m")

	uiWithoutWeekends := &UI{
		Calendar: &cal.BusinessCalendar{
			WorkdayFunc:      WorkdayOnlyWeekdays,
			WorkdayStartFunc: WorkdayStart,
			WorkdayEndFunc:   WorkdayEnd,
		},
	}
	st.Assert(t, uiWithoutWeekends.getTimeToFirstReview("Batman", "", false, timelineItems, reviews), "15h11m")
}

func Test_getTimeToFirstReview_Draft(t *testing.T) {
	var timelineItems = TimelineItems{
		Nodes: TimelineItemNodes{},
	}
	var reviews = Reviews{
		Nodes: ReviewNodes{
			{
				Author:    Author{Login: "Joker"},
				CreatedAt: "2022-03-22T15:11:09Z",
				State:     "COMMENTED",
			},
		},
	}

	ui := &UI{}
	st.Assert(t, ui.getTimeToFirstReview("Batman", "", true, timelineItems, reviews), "--")
}

func Test_getTimeToFirstReview_NoReviews(t *testing.T) {
	var timelineItems = TimelineItems{
		Nodes: TimelineItemNodes{
			{ReadyForReviewEvent{CreatedAt: "2022-03-21T15:11:09Z"}},
		},
	}
	var reviews = Reviews{
		Nodes: ReviewNodes{},
	}

	ui := &UI{}
	st.Assert(t, ui.getTimeToFirstReview("Batman", "", false, timelineItems, reviews), "--")
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

func Test_getFirstReviewToLastReview_OnlyAuthorReview(t *testing.T) {
	var reviews = Reviews{
		Nodes: ReviewNodes{
			{
				Author: Author{
					Login: "Batman",
				},
				CreatedAt: "2022-04-06T15:11:09Z",
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

	st.Assert(t, uiWithWeekends.getFirstReviewToLastReview("Batman", reviews), DefaultEmptyCell)
}

func Test_getFirstReviewToLastReview_NoApprovals(t *testing.T) {
	var reviews = Reviews{
		Nodes: ReviewNodes{
			{
				Author: Author{
					Login: "Joker",
				},
				CreatedAt: "2022-04-06T15:11:09Z",
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

	st.Assert(t, uiWithWeekends.getFirstReviewToLastReview("Batman", reviews), DefaultEmptyCell)
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
