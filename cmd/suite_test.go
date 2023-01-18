package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gopkg.in/h2non/gock.v1"
)

const (
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

func gqlSearchQueryMatcher(owner, repo, start, end string) func(req *http.Request, ereq *gock.Request) (bool, error) {
	return func(req *http.Request, ereq *gock.Request) (bool, error) {
		var gqlRequest GQLRequest

		var body, err = io.ReadAll(req.Body)
		err = json.Unmarshal(body, &gqlRequest)

		return gqlRequest.Variables.Query == fmt.Sprintf("repo:%s/%s type:pr merged:%s..%s",
			owner,
			repo,
			start,
			end), err
	}
}
